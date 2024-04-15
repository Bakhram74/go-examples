package repo

import (
	"errors"
	"fmt"
	"single-window/internal/entity"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

const (
	disputesTableName      = "disputes.dispute"
	disputesRolesTableName = "disputes.dispute_role"
	disputesChatsTableName = "disputes.chat"
	shortagesTableName     = "shortages.shortage"
)

type DisputeRepo struct {
	db *gorm.DB
}

func NewDisputeRepo(db *gorm.DB) *DisputeRepo {
	return &DisputeRepo{db}
}

func (d *DisputeRepo) GetDisputesV1(claims entity.Claims, limit, offset int, status entity.EntityStatus) ([]entity.DisputeList, error) {
	var disputes []entity.DisputeList
	baseQuery := `SELECT d.dispute_id,
					u.fullname AS complainant_name,
					u.employee_id AS complainant_employee_id,
					s.goods_id,
					s.tare_id,
					s.tare_type,
					s.currency_code,
					d.created_at,
					d.closed_at,
					l.lostreason_val,
					s.lost_amount,
					d.status,
					us.fullname AS responsible_person_name,
					o.organization_title
				FROM disputes.dispute d
				LEFT JOIN disputes.dispute_role dr ON d.dispute_id = dr.dispute_id AND dr.dispute_role = 'complainant'
				LEFT JOIN subjects.user u ON dr.user_id = u.user_id
				LEFT JOIN shortages.shortage s ON d.shortage_id = s.shortage_id
				LEFT JOIN shortages.lostreason l ON s.lostreason_id = l.lostreason_id
				LEFT JOIN disputes.dispute_role dis ON d.dispute_id = dis.dispute_id AND dis.dispute_role = 'responsible_person'
				LEFT JOIN subjects.user us ON dis.user_id = us.user_id
				LEFT JOIN subjects.organization o ON o.organization_id = d.organization_id`

	var query string
	var args []interface{}

	// TODO: переделать конструирование строки под strings.Builder
	if claims.UserRole == "arbitr" || claims.UserRole == "responsible_person" {
		if status == "opened" {
			query = baseQuery + ` WHERE d.status = 'opened'`
		} else if status != "any" {
			query = baseQuery + " WHERE us.user_id = ? AND d.status = ?"
			args = append(args, claims.UserId, status)
		} else {
			query = baseQuery + " WHERE us.user_id = ? OR d.status = 'opened'"
			args = append(args, claims.UserId)
		}
	} else if claims.UserRole == "complainant" {
		query = baseQuery + ` WHERE u.user_id = ?`
		args = append(args, claims.UserId)

		if status != "any" {
			query += " AND d.status = ?"
			args = append(args, status)
		}
	}

	query += " ORDER BY d.created_at LIMIT ? OFFSET ?"
	args = append(args, limit, offset)

	err := d.db.Raw(query, args...).Find(&disputes).Error
	if err != nil {
		return nil, err
	}

	return disputes, nil
}

func (d *DisputeRepo) GetGuiltyWorkersV1(disputeId string) ([]string, error) {
	var guiltyWorkerNames []string

	baseQuery := fmt.Sprintf(`
	SELECT u.fullname as guilty_worker_names
	FROM disputes.dispute_role dr
	LEFT JOIN subjects.user u
		ON dr.dispute_role = 'guilty_worker'
			WHERE dr.dispute_id = '%s'
				AND dr.user_id = u.user_id;
	`, disputeId)

	err := d.db.Raw(baseQuery).Find(&guiltyWorkerNames).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get guilty workers names, err: %w", err)
	}

	return guiltyWorkerNames, nil
}

func (d *DisputeRepo) GetDisputeV1(disputeId *string, goodsId *int) (*entity.Dispute, error) {
	var dispute entity.Dispute

	baseQuery := fmt.Sprintln(`
	SELECT d.dispute_id,
		d.status,
		u.fullname AS complainant_name,
		ur.fullname AS responsible_person_name,
		usr_g.fullname AS guilty_responsible_person_name,
		d.is_shortage_canceled,
		d.created_at,
		d.closed_at,
		s.goods_id,
		s.tare_id,
		s.tare_type,
		s.lost_amount,
		s.currency_code,
		lr.lostreason_val
	FROM disputes.dispute d
	LEFT JOIN disputes.dispute_role drc ON drc.dispute_role = 'complainant' AND drc.dispute_id = d.dispute_id
	LEFT JOIN disputes.dispute_role drs	ON drs.dispute_role = 'responsible_person' AND drs.dispute_id = d.dispute_id
	LEFT JOIN disputes.dispute_role drg ON drg.dispute_role = 'guilty_responsible_person' AND drg.dispute_id = d.dispute_id
	LEFT JOIN subjects.user u ON u.user_id = drc.user_id
	LEFT JOIN subjects.user ur ON ur.user_id = drs.user_id
	LEFT JOIN subjects.user usr_g ON usr_g.user_id = drg.user_id
	LEFT JOIN shortages.shortage s ON s.shortage_id = d.shortage_id
	LEFT JOIN shortages.lostreason lr ON lr.lostreason_id = s.lostreason_id
	`)

	groupQuery := `
	GROUP BY d.dispute_id,
			u.fullname,
			ur.fullname,
			usr_g.fullname,
			s.goods_id,
			s.tare_id,
			s.tare_type,
			s.lost_amount,
			s.currency_code,
			lr.lostreason_val;
	`

	var query string

	if disputeId != nil {
		query = baseQuery + fmt.Sprintf(`WHERE d.dispute_id = '%s'`, *disputeId) + groupQuery
	} else if goodsId != nil {
		query = baseQuery + fmt.Sprintf(`WHERE s.goods_id = '%d'`, *goodsId) + groupQuery
	}

	err := d.db.Raw(query).Find(&dispute).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get dispute, err: %w", err)
	}

	return &dispute, nil
}

func (d *DisputeRepo) CreateDisputeV1(createDisputeReqBody *entity.CreateDisputeV1MultipartBody, userId string, disputeAttachmentPath *string) (*entity.CreateDispute, error) {
	disputeId := uuid.NewString()
	disputeCreatedAt := time.Now()
	disputeToCreate := entity.CreateDispute{
		DisputeId:          disputeId,
		CreatedAt:          disputeCreatedAt,
		OrganizationId:     createDisputeReqBody.OrganizationId,
		IsShortageCanceled: false,
		IsArbitrInvited:    false,
		IsDisputeReopened:  false,
		ShortageId:         createDisputeReqBody.ShortageId,
		Status:             entity.Opened,
		ReopenedAt:         nil,
		ClosedAt:           nil,
	}

	shortageToFind := entity.Shortage{}
	err := d.db.Table(shortagesTableName).First(&shortageToFind, "shortage_id = ? AND user_id = ?", disputeToCreate.ShortageId, userId).Error
	if err != nil {
		return nil, fmt.Errorf("failed to find shortage with id = %s on user_id = %s", disputeToCreate.ShortageId, userId)
	}

	disputeToFind := entity.CreateDispute{}
	err = d.db.Table(disputesTableName).First(&disputeToFind, "shortage_id = ?", disputeToCreate.ShortageId).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("failed to get dispute by shortage_id, err: %w", err)
	}

	if disputeToFind.DisputeId != "" {
		return nil, fmt.Errorf("dispute on shartage_id: %s has already been created", disputeToFind.ShortageId)
	}

	err = d.db.Transaction(func(tx *gorm.DB) error {
		err = tx.Table(disputesTableName).Create(&disputeToCreate).Error
		if err != nil {
			return fmt.Errorf("failed to create dispute, error: %w", err)
		}

		disputeRole := entity.DisputeRole{
			UserId:      userId,
			DisputeRole: entity.Complainant,
			CreatedAt:   disputeCreatedAt,
			DisputeId:   disputeId,
		}

		err = tx.Table(disputesRolesTableName).Create(disputeRole).Error
		if err != nil {
			return fmt.Errorf("failed to create dispute role, err: %w", err)
		}

		disputeChat := entity.DisputeChat{
			MessageId:      uuid.NewString(),
			DisputeId:      disputeId,
			SenderId:       userId,
			AttachmentPath: disputeAttachmentPath,
			MessageBody:    &createDisputeReqBody.MessageBody,
			CreatedAt:      disputeCreatedAt,
		}

		err = tx.Table(disputesChatsTableName).Create(&disputeChat).Error
		if err != nil {
			return fmt.Errorf("failed to create dispute chat, error: %w", err)
		}

		err = tx.Table(shortagesTableName).Model(&entity.Shortage{}).Where("shortage_id = ?", disputeToCreate.ShortageId).Update("is_disputed", true).Error
		if err != nil {
			return fmt.Errorf("failed to modify shortage 'is_disputed' property, err: %w", err)
		}

		return nil
	})

	return &disputeToCreate, err
}

func (d *DisputeRepo) DeleteDisputeV1(dispute *entity.CreateDispute) (*entity.CreateDispute, error) {
	err := d.db.Table(disputesTableName).Delete(dispute, "dispute_id = ?", dispute.DisputeId).Error
	if err != nil {
		return nil, fmt.Errorf("failed to delete dispute, err: %w", err)
	}

	return dispute, nil
}

func (d *DisputeRepo) CloseDisputeV1(disputeId string, userIds []string) (*entity.CloseDispute, error) {
	currentTime := time.Now()

	var dispute *entity.CloseDispute
	var guiltyWorkres []entity.DisputeRole

	for _, userID := range userIds {
		guiltyWorkres = append(guiltyWorkres, entity.DisputeRole{
			UserId:      userID,
			CreatedAt:   currentTime,
			DisputeRole: "guilty_worker",
			DisputeId:   disputeId,
		})
	}

	err := d.db.Transaction(func(tx *gorm.DB) error {
		err := tx.Table(disputesRolesTableName).
			Where("dispute_role = 'guilty_worker'").
			Where("dispute_id = ?", disputeId).
			Delete(&entity.DisputeRole{}).
			Error
		if err != nil {
			return fmt.Errorf("failed delete guilty dispute role, err: %w", err)
		}

		if len(guiltyWorkres) > 0 {
			err = tx.Table(disputesRolesTableName).Create(&guiltyWorkres).Error
			if err != nil {
				return fmt.Errorf("failed to create guilty dispute role, err: %w", err)
			}
		} else {
			err = tx.Table(disputesTableName).Where("dispute_id = ?", disputeId).Update("is_shortage_canceled", true).Error
			if err != nil {
				return fmt.Errorf("failed change is_shortage_canceled, err: %w", err)
			}
		}

		err = tx.Table(disputesTableName).
			Where("dispute_id = ?", disputeId).
			Update("status", "closed").
			Update("closed_at", currentTime).
			Error
		if err != nil {
			return fmt.Errorf("failed to update dispute, err: %w", err)
		}

		err = tx.Table("disputes.dispute AS d").
			Select("d.*").
			Where("dispute_id = ?", disputeId).
			First(&dispute).
			Error
		if err != nil {
			return fmt.Errorf("failed to get dispute, err: %w", err)
		}

		err = tx.Table("disputes.dispute_role as dr").
			Select("u.fullname").
			Joins("LEFT JOIN subjects.user AS u ON dr.user_id = u.user_id").
			Where("dr.dispute_id = ?", disputeId).
			Where("dr.dispute_role = 'guilty_worker'").
			Find(&dispute.GuiltyWorkerNames).
			Error
		if err != nil {
			return fmt.Errorf("failed to get guilty_worker_names, err: %w", err)
		}

		return nil
	})

	return dispute, err
}

func (d *DisputeRepo) GetShortagesV1(userId string, limits, offset int) ([]entity.Shortage, error) {
	var shortages []entity.Shortage

	query := `SELECT s.shortage_id,
					 s.user_id,
					 s.goods_id,
					 s.tare_id,
					 s.tare_type,
					 s.lostreason_id,
					 s.currency_code,
					 s.lost_amount,
					 s.is_disputed,
					 s.created_at
			  FROM shortages.shortage s
			  WHERE s.user_id = ?
			  AND s.is_disputed = false
			  LIMIT ? OFFSET ?`

	err := d.db.Raw(query, userId, limits, offset).Find(&shortages).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get shortages, err: %w", err)
	}

	return shortages, nil
}
