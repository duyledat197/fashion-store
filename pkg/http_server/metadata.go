package http_server

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"google.golang.org/grpc/metadata"

	"trintech/review/pkg/http_server/xcontext"
	stringutil "trintech/review/pkg/string_util"
	"trintech/review/pkg/token_util"
)

const (
	AUTHORIZATION = "Authorization"
	BEARER        = "Bearer"
)

const (
	MDUserIDKey     = "user_id"
	MDIpKey         = "ip"
	MDUserAgent     = "user-agent"
	MDRoleKey       = "role"
	MDXForwardedFor = "x-forwarded-for"
)

// DataResponse ...
func DataResponse(w http.ResponseWriter, data any) {
	resp := &Response{
		Data: data,
	}

	jData, _ := json.Marshal(resp)

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(jData)
}

type mapMetaDataFunc func(context.Context, *http.Request) metadata.MD

// MapMetaDataWithBearerToken ...
func MapMetaDataWithBearerToken(authenticator token_util.Authenticator[*xcontext.UserInfo]) mapMetaDataFunc {
	return func(ctx context.Context, r *http.Request) metadata.MD {
		md, _ := metadata.FromIncomingContext(ctx)

		md = metadata.Join(md, ImportSessionToMD(&xcontext.Session{
			IP:        stringutil.Coalesce(md.Get(MDXForwardedFor)...),
			UserAgent: stringutil.Coalesce(md.Get(MDUserAgent)...),
		}))

		authorization := r.Header.Get(AUTHORIZATION)

		if authorization != "" {
			schema, token, isValid := strings.Cut(authorization, " ")
			if schema == BEARER && isValid {
				payload, err := authenticator.Verify(token)
				if err == nil {
					md = metadata.Join(md, ImportUserInfoToMD(payload))
				}
			}
		}

		return md
	}
}

// ImportUserInfoToMD ...
func ImportUserInfoToMD(payload *xcontext.UserInfo) metadata.MD {
	md := metadata.Pairs(
		MDUserIDKey, fmt.Sprint(payload.UserID), // append userID
		MDRoleKey, payload.Role, // append role
	)

	return md
}

// ImportSessionToMD ...
func ImportSessionToMD(payload *xcontext.Session) metadata.MD {
	md := metadata.Pairs(
		MDIpKey, string(payload.IP), // append userID
		MDUserAgent, payload.UserAgent, // append user agent
	)

	return md
}
