package entity

import (
	"single-window/pkg/httpclient"
	"strconv"
)

func (headers AuthUserParams) ConvertHeadersToStringMap() map[string]string {
	const (
		xRealIPHeader        = "X-Real-IP"
		deviceIDHeader       = "deviceId"
		xRequestIDHeader     = "X-Request-ID"
		deviceTokenHeader    = "deviceToken"
		deviceNameHeader     = "deviceName"
		wbApptypeHeader      = "wb-apptype"
		wbAppversionHeader   = "wb-appversion"
		xNoSessionHeader     = "X-No-Session"
		xForwardedHostHeader = "X-Forwarded-Host,omitempty"
		xForwardedPathHeader = "X-Forwarded-Path"
		acceptHeader         = "accept"
	)

	requestHeaders := make(map[string]string)
	requestHeaders[acceptHeader] = string(httpclient.ApplicationJson)
	requestHeaders[xRealIPHeader] = headers.XRealIP
	if headers.DeviceId != nil {
		requestHeaders[deviceIDHeader] = *headers.DeviceId
	}
	if headers.DeviceToken != nil {
		requestHeaders[deviceTokenHeader] = *headers.DeviceToken
	}

	if headers.DeviceName != nil {
		requestHeaders[deviceNameHeader] = *headers.DeviceName
	}

	if headers.WbApptype != nil {
		requestHeaders[wbApptypeHeader] = *headers.WbApptype
	}

	if headers.WbAppversion != nil {
		requestHeaders[wbAppversionHeader] = *headers.WbAppversion
	}

	if headers.XNoSession != nil {
		requestHeaders[xNoSessionHeader] = strconv.FormatBool(*headers.XNoSession)
	}

	if headers.XForwardedHost != nil {
		requestHeaders[xForwardedHostHeader] = *headers.XForwardedHost
	}

	if headers.XForwardedPath != nil {
		requestHeaders[xForwardedPathHeader] = *headers.XForwardedPath
	}
	return requestHeaders
}
