-- Вставка тестовых данных в схему disputes
-- Вставка тестовых данных споров
INSERT INTO disputes.dispute
VALUES  -- Вставка новых споров со статусом 'opened'
       ('9cfcfb2a-9e4b-44c9-b1f3-337a97a00e12',
        '48142fe2-42d7-40ed-84b5-05c5eac09382',
        '733378f1-77fd-4460-bf24-c14f6687c2fb',
        FALSE,
        FALSE,
        'opened',
        '2023-11-10 17:09:14.00000',
        NULL,
        NULL,
        FALSE),
       ('75402413-d3e1-4723-ba26-bec4a8b89d9a',
        '8a53b35e-51ee-4c29-80d7-76acafc9b027',
        'c721bfdf-f048-403b-a846-83e4977410ac',
        FALSE,
        FALSE,
        'opened',
        '2023-11-10 17:09:14.00000',
        NULL,
        NULL,
        FALSE),
        -- Вставка споров, взятых в работу, со статусом 'in_work'
       ('91451770-423e-42ca-ab17-7f05b8c938fe',
        '743979cd-e7c2-4b01-a847-8b93c0962013',
        '733378f1-77fd-4460-bf24-c14f6687c2fb',
        FALSE,
        FALSE,
        'in_work',
        '2023-11-10 17:09:14.00000',
        NULL,
        NULL,
        FALSE),
       ('ec90fdca-19cc-4645-aca1-c895d89f4e04',
        '85ae15ea-eb0b-4fcb-99d2-e7852c64fa0a',
        '84b9a947-d4ab-4f70-889b-57ad6d8c7db5',
        FALSE,
        FALSE,
        'in_work',
        '2023-11-10 17:09:14.00000',
        NULL,
        NULL,
        FALSE),
        -- Вставка споров, взятых в работу, с открытым приглашением арбитра
       ('3b15aa77-f59e-4c31-b16d-a44e25396052',
        '627cdce8-eddf-456e-a646-eade0060e549',
        '825d29fb-5dc9-4c73-ad48-bc24a42cbe0e',
        FALSE,
        FALSE,
        'in_work',
        '2023-11-10 17:09:14.00000',
        NULL,
        NULL,
        TRUE),
        -- Вставка споров, взятых в работу, со взятым в работу приглашением арбитра
       ('2231a2bf-c771-4909-8611-8ef977ca1ef6',
        '5f26bb19-3995-40cc-9d95-a0aea566d69b',
        '7bb392f4-1e87-447e-9286-96651e8a1369',
        FALSE,
        FALSE,
        'in_work',
        '2023-11-13 10:09:14.00000',
        NULL,
        NULL,
        TRUE),
        -- Вставка завершенных споров с приглашением арбитра и анулированием списания
       ('cea09d66-f41a-4b94-a955-0d33610153e1',
        '7170c8dd-ffb4-47e9-8dd2-e6ffee9beec6',
        '7bb392f4-1e87-447e-9286-96651e8a1369',
        TRUE,
        FALSE,
        'closed',
        '2023-11-09 10:09:14.00000',
        '2023-11-13 11:03:14.00000',
        NULL,
        TRUE),
        -- Вставка завершенных споров с приглашением арбитра и назначением виновного лица
       ('cbc79f25-ab55-4fc6-b149-13d63ef698d0',
        'a37d2349-d207-40c2-9264-c7a1e3c99754',
        '8efe12d9-62e6-44fe-91d0-bfaac764433f',
        FALSE,
        FALSE,
        'closed',
        '2023-11-09 10:09:14.00000',
        '2023-11-13 11:03:14.00000',
        NULL,
        TRUE),
        -- Вставка завершенных споров с приглашением арбитра и назначением виновного лица + назначение виновного ответственного лица
       ('4ee79d3f-04a7-4baa-903d-cf43de15b654',
        '60ff4317-058d-4102-84ec-0dc148c0cfea',
        '733378f1-77fd-4460-bf24-c14f6687c2fb',
        FALSE,
        FALSE,
        'closed',
        '2023-11-09 10:09:14.00000',
        '2023-11-13 11:03:14.00000',
        NULL,
        TRUE),
        -- Вставка завершенных споров с назначением виновного лица
       ('b7e4fb40-bb27-413e-b015-87f6ce0eae6d',
        'ef40f67e-c819-4c36-ad1a-1fc0b8096add',
        '733378f1-77fd-4460-bf24-c14f6687c2fb',
        FALSE,
        FALSE,
        'closed',
        '2023-11-09 10:09:14.00000',
        '2023-11-13 11:03:14.00000',
        NULL,
        FALSE),
        -- Вставка переоткрытых споров в статусе "открыт"
       ('9e41266a-ff8f-4a3f-8b4d-05afbafdfe2a',
        'b768c8e6-22f0-4fdf-843e-30a4b2c94ef5',
        '733378f1-77fd-4460-bf24-c14f6687c2fb',
        FALSE,
        TRUE,
        'opened',
        '2023-11-09 10:09:14.00000',
        NULL,
        '2023-11-13 11:03:14.00000',
        FALSE),
        -- Вставка переоткрытых споров в статусе "в работе"
       ('514fc3ba-f498-46bc-9e07-7e6b2e1178fb',
        'e6d359cb-e097-4899-b1e6-da1b58ccad31',
        '733378f1-77fd-4460-bf24-c14f6687c2fb',
        FALSE,
        TRUE,
        'in_work',
        '2023-11-09 10:09:14.00000',
        NULL,
        '2023-11-13 11:03:14.00000',
        FALSE),
        -- Вставка переоткрытого спора с анулированным списанием в статусе "закрыт"
       ('b5eadb23-c008-4d53-b2ba-8d1e86c64906',
        '81fdc1b0-f739-4df1-bab0-22c4bd4acb5f',
        '733378f1-77fd-4460-bf24-c14f6687c2fb',
        TRUE,
        TRUE,
        'closed',
        '2023-11-09 10:09:14.00000',
        '2023-11-15 14:18:14.00000',
        '2023-11-13 11:03:14.00000',
        FALSE)
ON CONFLICT DO NOTHING;

-- Вставка тестовых записей ролей пользователей в спорах
INSERT INTO disputes.dispute_role
VALUES ('88e75a84-956a-4b25-b991-c9fd0e848715',
        '9cfcfb2a-9e4b-44c9-b1f3-337a97a00e12',
        'complainant',
        '2023-11-10 17:09:14.00000'),
       ('0d7a4fd9-b0be-45f9-91b5-bb1b95a601f3',
        '75402413-d3e1-4723-ba26-bec4a8b89d9a',
        'complainant',
        '2023-11-10 17:09:14.00000'),
       ('88e75a84-956a-4b25-b991-c9fd0e848715',
        '91451770-423e-42ca-ab17-7f05b8c938fe',
        'complainant',
        '2023-11-10 17:09:14.00000'),
       ('a9d1e540-4a4e-4d55-883f-6cda83336d45',
        '91451770-423e-42ca-ab17-7f05b8c938fe',
        'responsible_person',
        '2023-11-10 18:09:14.00000'),
       ('0f7d2818-bdbf-440c-b3dd-7a0009888e64',
        'ec90fdca-19cc-4645-aca1-c895d89f4e04',
        'complainant',
        '2023-11-10 17:09:14.00000'),
       ('a9d1e540-4a4e-4d55-883f-6cda83336d45',
        'ec90fdca-19cc-4645-aca1-c895d89f4e04',
        'responsible_person',
        '2023-11-10 18:36:14.00000'),
       ('9004a9ec-a656-4ff9-aa47-c3a7af658a29',
        '3b15aa77-f59e-4c31-b16d-a44e25396052',
        'complainant',
        '2023-11-10 17:09:14.00000'),
       ('c90fed91-414f-49bc-b3d4-f248fd776901',
        '3b15aa77-f59e-4c31-b16d-a44e25396052',
        'responsible_person',
        '2023-11-10 17:48:14.00000'),
       ('d4c22a44-2c8d-4dab-a01f-76da538c8efe',
        '2231a2bf-c771-4909-8611-8ef977ca1ef6',
        'complainant',
        '2023-11-13 10:09:14.00000'),
       ('09eb21a0-123d-4eab-847a-0795f5b97e64',
        '2231a2bf-c771-4909-8611-8ef977ca1ef6',
        'responsible_person',
        '2023-11-13 13:09:14.00000'),
       ('b8fb4536-6a6f-47e6-86f5-de9300c32768',
        'cea09d66-f41a-4b94-a955-0d33610153e1',
        'complainant',
        '2023-11-09 10:09:14.00000'),
       ('4cb24169-cb16-4bf8-a4fd-0fa3cef4fbad',
        'cea09d66-f41a-4b94-a955-0d33610153e1',
        'responsible_person',
        '2023-11-10 20:09:14.00000'),
       ('f7fb6a8d-aee3-4ffb-ab4f-d3a891ea444c',
        'cbc79f25-ab55-4fc6-b149-13d63ef698d0',
        'complainant',
        '2023-11-09 10:09:14.00000'),
       ('a218592a-86f6-45ac-9e9d-63a7b70c2a0b',
        'cbc79f25-ab55-4fc6-b149-13d63ef698d0',
        'responsible_person',
        '2023-11-10 18:24:14.00000'),
       ('f7fb6a8d-aee3-4ffb-ab4f-d3a891ea444c',
        'cbc79f25-ab55-4fc6-b149-13d63ef698d0',
        'guilty_worker',
        '2023-11-13 11:03:14.00000'),
       ('0f7d2818-bdbf-440c-b3dd-7a0009888e64',
        '4ee79d3f-04a7-4baa-903d-cf43de15b654',
        'complainant',
        '2023-11-09 10:09:14.00000'),
       ('a9d1e540-4a4e-4d55-883f-6cda83336d45',
        '4ee79d3f-04a7-4baa-903d-cf43de15b654',
        'responsible_person',
        '2023-11-09 18:09:14.00000'),
       ('0f7d2818-bdbf-440c-b3dd-7a0009888e64',
        '4ee79d3f-04a7-4baa-903d-cf43de15b654',
        'guilty_worker',
        '2023-11-13 11:03:14.00000'),
       ('a9d1e540-4a4e-4d55-883f-6cda83336d45',
        '4ee79d3f-04a7-4baa-903d-cf43de15b654',
        'guilty_responsible_person',
        '2023-11-13 11:03:14.00000'),
       ('88e75a84-956a-4b25-b991-c9fd0e848715',
        'b7e4fb40-bb27-413e-b015-87f6ce0eae6d',
        'complainant',
        '2023-11-09 10:09:14.00000'),
       ('a9d1e540-4a4e-4d55-883f-6cda83336d45',
        'b7e4fb40-bb27-413e-b015-87f6ce0eae6d',
        'responsible_person',
        '2023-11-11 11:42:14.00000'),
       ('88e75a84-956a-4b25-b991-c9fd0e848715',
        'b7e4fb40-bb27-413e-b015-87f6ce0eae6d',
        'guilty_worker',
        '2023-11-13 11:03:14.00000'),
       ('88e75a84-956a-4b25-b991-c9fd0e848715',
        '9e41266a-ff8f-4a3f-8b4d-05afbafdfe2a',
        'complainant',
        '2023-11-09 10:09:14.00000'),
       ('88e75a84-956a-4b25-b991-c9fd0e848715',
        '514fc3ba-f498-46bc-9e07-7e6b2e1178fb',
        'complainant',
        '2023-11-09 10:09:14.00000'),
       ('adaae422-3686-4605-81d7-d8c6471896f2',
        '514fc3ba-f498-46bc-9e07-7e6b2e1178fb',
        'responsible_person',
        '2023-11-13 11:03:14.00000'),
       ('88e75a84-956a-4b25-b991-c9fd0e848715',
        'b5eadb23-c008-4d53-b2ba-8d1e86c64906',
        'complainant',
        '2023-11-09 10:09:14.00000'),
       ('adaae422-3686-4605-81d7-d8c6471896f2',
        'b5eadb23-c008-4d53-b2ba-8d1e86c64906',
        'responsible_person',
        '2023-11-13 11:03:14.00000')
ON CONFLICT DO NOTHING;

-- Вставка тестовых записей приглашений арбитров
INSERT INTO disputes.arbitrinvitation
VALUES ('7972c7a1-14ba-49fa-9666-9e8a97cdb87c',
        '825d29fb-5dc9-4c73-ad48-bc24a42cbe0e',
        '3b15aa77-f59e-4c31-b16d-a44e25396052',
        'opened',
        '2023-11-10 18:40:00.00000',
        NULL,
        NULL,
        NULL),
       ('60a164c2-75d0-4cba-ad64-c5d0866535b3',
        '7bb392f4-1e87-447e-9286-96651e8a1369',
        '2231a2bf-c771-4909-8611-8ef977ca1ef6',
        'in_work',
        '2023-11-12 10:40:00.00000',
        '2023-11-13 10:45:00.00000',
        NULL,
        'adaae422-3686-4605-81d7-d8c6471896f2'),
       ('a9429f3c-2ed6-479a-87dd-14abf51ba20e',
        '7bb392f4-1e87-447e-9286-96651e8a1369',
        'cea09d66-f41a-4b94-a955-0d33610153e1',
        'closed',
        '2023-11-11 10:40:00.00000',
        '2023-11-12 09:03:14.00000',
        '2023-11-13 11:03:14.00000',
        'adaae422-3686-4605-81d7-d8c6471896f2'),
       ('06391d8d-2d61-4bcc-89a0-47e90d977deb',
        '8efe12d9-62e6-44fe-91d0-bfaac764433f',
        'cbc79f25-ab55-4fc6-b149-13d63ef698d0',
        'closed',
        '2023-11-11 10:40:00.00000',
        '2023-11-12 09:03:14.00000',
        '2023-11-13 11:03:14.00000',
        'adaae422-3686-4605-81d7-d8c6471896f2'),
       ('5e1410ed-3850-4dd0-bc2e-b72a1a9e99ed',
        '733378f1-77fd-4460-bf24-c14f6687c2fb',
        '4ee79d3f-04a7-4baa-903d-cf43de15b654',
        'closed',
        '2023-11-11 10:40:00.00000',
        '2023-11-12 09:03:14.00000',
        '2023-11-13 11:03:14.00000',
        'adaae422-3686-4605-81d7-d8c6471896f2')
ON CONFLICT DO NOTHING;

-- Вставка тестовых записей чата
INSERT INTO disputes.chat
VALUES ('547e7da3-bd89-40a5-902d-6fe4a9b03949',
        '91451770-423e-42ca-ab17-7f05b8c938fe',
        '88e75a84-956a-4b25-b991-c9fd0e848715',
        '2023-11-10 18:10:25.00120',
        'Добрый вечер! Я не согласен со списанием, потому что в этот день меня не было на работе - я болел! Больничный листок ест в кадрах.',
        NULL),
       ('9d6075b9-f82a-4b79-b687-73a6b926e77d',
        '91451770-423e-42ca-ab17-7f05b8c938fe',
        'a9d1e540-4a4e-4d55-883f-6cda83336d45',
        '2023-11-10 19:25:25.00120',
        'Добрый вечер! Вас понял, сейчас запрошу в кадрах. Ждите.',
        NULL),
       ('1be705a1-7768-4e05-9bac-5f69569e3ca3',
        '91451770-423e-42ca-ab17-7f05b8c938fe',
        '88e75a84-956a-4b25-b991-c9fd0e848715',
        '2023-11-11 09:36:21.00120',
        NULL,
        'mock/attachments/chat/91451770-423e-42ca-ab17-7f05b8c938fe_111123093625.jpeg'),
       ('cdce5fab-1d32-43f6-9651-cd9b0f38358a',
        '91451770-423e-42ca-ab17-7f05b8c938fe',
        'a9d1e540-4a4e-4d55-883f-6cda83336d45',
        '2023-11-11 09:36:25.00120',
        'Здравствуйте! Отлично. Спасибо за предоставленные документы!',
        NULL),

       ('a7d739f9-bbc2-4f58-b43a-2c97ed59da2e',
        'ec90fdca-19cc-4645-aca1-c895d89f4e04',
        '0f7d2818-bdbf-440c-b3dd-7a0009888e64',
        '2023-11-10 17:25:25.00120',
        'Добрый день! Я не согласен со списанием, так как в день обработки списанного шк я работал на другом блоке.',
        NULL),
       ('944e16a6-ebbf-459c-a143-aace391e0d0a',
        'ec90fdca-19cc-4645-aca1-c895d89f4e04',
        'a9d1e540-4a4e-4d55-883f-6cda83336d45',
        '2023-11-10 17:28:25.00120',
        'Здравствуйте! Предоставьте, пожалуйста, какие-нибудь доказательства.',
        NULL),

       ('9d19ab92-9684-4ff0-a390-d2fc871ed604',
        '3b15aa77-f59e-4c31-b16d-a44e25396052',
        '9004a9ec-a656-4ff9-aa47-c3a7af658a29',
        '2023-11-10 19:42:25.00120',
        'Добрый вечер! Я не согласен со списанием, потому что я обработал шк и тара засветилась на другом офисе (шк не был утерян на моем этапе обработки).',
        NULL)
ON CONFLICT DO NOTHING;

-- Вставка тестовых записей запросов
INSERT INTO disputes.revision
VALUES -- Вставка не взятых в работу запросов
       ('ae5aa9b8-2458-46df-bd6e-0696abae00bc',
        'c721bfdf-f048-403b-a846-83e4977410ac',
        '91451770-423e-42ca-ab17-7f05b8c938fe',
        'opened',
        '2023-11-10 19:33:25.00120',
        NULL,
        NULL,
        NULL),
        ('a274bbac-b552-46ed-8773-a9a8633d2292',
         '733378f1-77fd-4460-bf24-c14f6687c2fb',
         'ec90fdca-19cc-4645-aca1-c895d89f4e04',
         'opened',
         '2023-11-10 19:33:25.00120',
        NULL,
        NULL,
        NULL),
       -- Вставка взятых в работу запросов
       ('3c2f0ff1-f6f2-4aa2-980c-4bff7037ad39',
        '8efe12d9-62e6-44fe-91d0-bfaac764433f',
        'ec90fdca-19cc-4645-aca1-c895d89f4e04',
        'in_work',
        '2023-11-10 20:48:25.00120',
        '2023-11-11 08:33:25.00120',
        NULL,
        'a218592a-86f6-45ac-9e9d-63a7b70c2a0b'),
        ('3d918e18-9bf2-426a-9035-0d5c70a41974',
         '733378f1-77fd-4460-bf24-c14f6687c2fb',
         'ec90fdca-19cc-4645-aca1-c895d89f4e04',
         'in_work',
         '2023-11-10 20:48:25.00120',
         '2023-11-11 06:33:25.00120',
         NULL,
         'a9d1e540-4a4e-4d55-883f-6cda83336d45'),
       -- Вставка завершенных запросов
       ('96abe567-c77d-44b0-82d0-de5ca83194e8',
        '84b9a947-d4ab-4f70-889b-57ad6d8c7db5',
        '3b15aa77-f59e-4c31-b16d-a44e25396052',
        'closed',
        '2023-11-10 20:48:25.00120',
        '2023-11-11 08:33:25.00120',
        '2023-11-12 19:33:25.00120',
        '4cb24169-cb16-4bf8-a4fd-0fa3cef4fbad')
ON CONFLICT DO NOTHING;

-- Вставка тестовых записей корреспонденции в рамках запросов
INSERT INTO disputes.correspondence
VALUES ('0e9884b9-b034-4b74-9123-d957a7b10bd5',
        'ae5aa9b8-2458-46df-bd6e-0696abae00bc',
        'a9d1e540-4a4e-4d55-883f-6cda83336d45',
        '2023-11-10 19:33:25.00120',
        'Добрый вечер! Предоставьте запись с камеры видеонаблюдения на блоке 123 офиса 123654 МХ 987654321 в период 11:00-11:10 07.09.2023.',
        NULL),

       ('6be6d3ce-99e1-4019-a7ce-68978c63cc54',
        '3c2f0ff1-f6f2-4aa2-980c-4bff7037ad39',
        'a9d1e540-4a4e-4d55-883f-6cda83336d45',
        '2023-11-11 08:21:25.00120',
        'Здравствуйте! Предоставьте фото\видео с камер видеонаблюдения на блоке 321 офиса 96325 МХ 365896325  в период 21:00-21:10 03.09.2023.',
        NULL),
       ('716d4719-a15c-4bb4-91ce-d6d1eab2a472',
        '3c2f0ff1-f6f2-4aa2-980c-4bff7037ad39',
        'a218592a-86f6-45ac-9e9d-63a7b70c2a0b',
        '2023-11-11 11:21:25.00120',
        'Доброе утро! Предоставляю видео, так как фото у нас нет.',
        'mock/attachments/revision/3c2f0ff1-f6f2-4aa2-980c-4bff7037ad39_111123112125.mp4'),
       ('794956af-5ed5-43fa-9d28-8087b012dac2',
        '3c2f0ff1-f6f2-4aa2-980c-4bff7037ad39',
        'a9d1e540-4a4e-4d55-883f-6cda83336d45',
        '2023-11-11 15:36:25.00120',
        'Спасибо! Можно еще видео в период 21:10-21:15 03.09.2023?',
        NULL),
       ('16f9e4bb-026d-4b04-9691-62234998f319',
        '3c2f0ff1-f6f2-4aa2-980c-4bff7037ad39',
        'a218592a-86f6-45ac-9e9d-63a7b70c2a0b',
        '2023-11-11 18:48:25.00120',
        'Конечно, прикрепляю.',
        'mock/attachments/revision/3c2f0ff1-f6f2-4aa2-980c-4bff7037ad39_111123184825.mp4')
ON CONFLICT DO NOTHING;