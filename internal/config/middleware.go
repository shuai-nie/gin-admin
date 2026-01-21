package config

type Middleware struct {
	Recovery struct {
		Skip int
	}
	CORS struct {
		Enable                 bool
		AllowAllOrigins        bool
		AllowOrigins           []string
		AllowMethods           []string
		AllowHeaders           []string
		AllowCredentials       bool
		ExposeHeaders          []string
		MaxAge                 int
		AllowWildcard          bool
		AllowBrowserExtensions bool
		AllowWebSockets        bool
		AllowFiles             bool
	}
	Trace struct {
		SkippedPathPrefixes []string
		RequestHeaderKey    string
		ResponseTraceKey    string
	}
	Logger struct {
		SkippedPathPrefixes      []string
		MaxOutputRequestBodyLen  int
		MaxOutputResponseBodyLen int
	}
	CopyBody struct {
		SkippedPathPrefixes []string
		MaxContentLen       int64
	}
	Auth struct {
		Disable             bool
		SkippedPathPrefixes []string
		SigningMethod       string
		SigningKey          string
		OldSigningKey       string
		Expired             int
		Store               struct {
			Type      string
			Delimiter string
			Memory    struct {
				CleanupInterval int
			}
			Badger struct {
				Path string
			}
			Redis struct {
				Addr     string
				Username string
				Password string
				DB       int
			}
		}
	}
	RateLimiter struct {
		Enable              bool
		SkippedPathPrefixes []string
		Period              int
		MaxRequestPerIP     int
		MaxRequestPerUser   int
		Store               struct {
			Type   string
			Memory struct {
				Expiration      int
				CleanupInterval int
			}
			Redis struct {
				Addr     string
				Username string
				Password string
				DB       int
			}
		}
	}
	Casbin struct {
		Disable             bool
		SkippedPathPrefixes []string
		LoadThread          int
		AutoLoadInterval    int
		ModelFile           string
		GenPolicyFile       string
	}

	Static struct {
		Dir string
	}
}
