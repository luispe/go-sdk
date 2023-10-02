package auth_test

import (
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/pomelo-la/go-toolkit/auth"
)

const (
	// The following private, public RSA key
	// and all tokens generated
	// random of https://dinochiesa.github.io/jwt/
	TestRSAPrivateKey = `
-----BEGIN PRIVATE KEY-----
MIIEvQIBADANBgkqhkiG9w0BAQEFAASCBKcwggSjAgEAAoIBAQCrq1LhfGO3cfr3
9Nb7NA4q+jbb9hFQWY5YByeDOUKvSGy6GBfnjMQiNS8Tn0V2GO3iRNH84LPsq9gZ
xxhANRWlX1idAIr4kSVSBP2NlwIlajwR2OT+mzSsateRunSTJFuxhmpFwtoBLpVl
/Vqf1Gv4fcVAlYVEQZMUR8F3Ye5UXGYQhMh61Iex6y0e+7jc0uOroZNwspcUlmms
gvwZSKrGae85kcIa7jxpUiJBYRrSi1yd9KNtdSfVeOvKKqWMRGv0lXmtPtl3fDzA
zA1COO0mX7lL+wipHZ6/ZEdzVsWxhYI/8njg403ZK2XfAKYr9dTVsyHNjyQ3MCFk
gpc8fLW3AgMBAAECggEABG+T364p769yL35FyvF6Lnc50TlzvU04ckaG1BJw0XUg
MFO3RffgpYy/0XD7lfMTLhYl09BoYHZuETDltoWvvoCH4Ylbf7c4bRlbPwNDIe9B
zzdci3EDeDNL5SRvCnz5BZTBBZLk0TXj8NO5AMJGFNfNzD7fDs0Da5W7+Cj4bs49
lEuxWJKdn5myuecmgydK3SOGhTrxJ1U+0Lmq0/sjk4eB+VqsrynK56kl8sqN9tDo
Ag8yo0l7oXf8upt+ykF37sW9+i1mReVN40ZlWEfn2UN3Gr6WMQrj+RGHdbhR/jHX
TLPVu70FLp5XC0agsaJGzJi58K3FuzxzU9k/pjiqQQKBgQDSYbSJdO26l/Ocjibe
yK8ufqWlEKcqC5BoG+B43GQE2kLfx9zRgOfykQuvFdRILprX/0FpmmVSHq0Ha8Fa
oefpif1BMzWLd6brnXt5tx9uhmFLq/poHrbFi+6IT1EO3Y9BrezvbVlky41H+ZTV
Rjo/EiyjmVeX5TKkoOH6R9tDYQKBgQDQ5LCnuk4SY59u1LWMC5GsehtGbWU2xgtx
UYpPDZCUdBwmegYhOCUMxTPqN9pxt5yCoFT1Jy3DKTCkmiyLtGvd590ort9cuBdC
o4Pyshj0T7VF3m60yVOyn/39Vc7DzR38ZN5IuweP++bueD2/OKEhu97dIO1lO5ss
qlhYXtKoFwKBgQDG3sNpeJXM4BzR7fJCgJRQsDlnKrHKZfoQ3+E2fqcxixzSKzzK
8j7QJlpUHJ95yExpSAqOh/ulQAgyTqMNSKVQNzembYD9IJMygMCa0wcsVG0euihQ
SlBdtyQ5yDiIg9oKrR2fSs/JHz2jPwN5BBTFUCnQUIDjvi48PzS+gTR8oQKBgH7X
ztkaTNvnuGEBMngmcj9sKfG67bGz0jDuFXDpSLiMRKesgtpbEExP1rVLUw6oMpYz
K0NtleEiutHIeHIgjTtC1s0kWqcfdahWSAHv2S1I1UbmyQxoD7WwZvcUyqekfqfK
zBsXzoDEsjZttvjNNzKXtL1LiDtnVVNq4JhQg9PjAoGAA8lEsRu4K7DEcHoGvoad
NAS3afjO2b3nIgJ52sLR21AK1jDf9gj4s5uQILH9ezpmqDgAt3ND9RJFfiL4w8+V
Thg093SKLCCb6UxkIv/QJ7L9/bptX3+oZyW+dWpcayICHAiXNgEeMP3eVVLYlCW9
NQ95+8yIxJ+tLCD2F9Nw2Tk=
-----END PRIVATE KEY-----
`
	TestRSAPublicKey = `
-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAq6tS4Xxjt3H69/TW+zQO
Kvo22/YRUFmOWAcngzlCr0hsuhgX54zEIjUvE59Fdhjt4kTR/OCz7KvYGccYQDUV
pV9YnQCK+JElUgT9jZcCJWo8Edjk/ps0rGrXkbp0kyRbsYZqRcLaAS6VZf1an9Rr
+H3FQJWFREGTFEfBd2HuVFxmEITIetSHsestHvu43NLjq6GTcLKXFJZprIL8GUiq
xmnvOZHCGu48aVIiQWEa0otcnfSjbXUn1XjryiqljERr9JV5rT7Zd3w8wMwNQjjt
Jl+5S/sIqR2ev2RHc1bFsYWCP/J44ONN2Stl3wCmK/XU1bMhzY8kNzAhZIKXPHy1
twIDAQAB
-----END PUBLIC KEY-----
`
)

func TestDecodeUserToken(t *testing.T) {
	type args struct {
		header string
		token  string
	}
	type mockClaims auth.Claims
	type expected struct {
		claims   *mockClaims
		ctxClaim string
		err      error
	}
	someUserClaims := mockClaims{Area: []string{"someBU"}, Email: "example@pomelo.la", Role: []string{"someRole"}}
	headerToken := "X-Auth-Token"
	userCtxClaim := "user"

	tests := []struct {
		name string
		args args
		want expected
	}{
		{
			name: "success user ctx audience",
			args: args{
				header: headerToken,
				token:  "eyJ0eXAiOiJKV1QiLCJhbGciOiJSUzI1NiJ9.eyJyb2xlIjpbInNvbWVSb2xlIl0sImVtYWlsIjoiZXhhbXBsZUBwb21lbG8ubGEiLCJhcmVhIjpbInNvbWVCVSJdLCJ0ZWFtcyI6WyJzb21lVGVhbSJdLCJleHAiOjM3MDA3MjI4NzUsImF1ZCI6InVzZXJjb250ZXh0YXBpIiwic2xhY2siOiJhYmMxMjMiLCJkZWxlZ2F0ZWQiOnRydWUsIm1hbmFnZXIiOm51bGwsImlkIjoiMTIzNGFjYmQiLCJzY29wZXMiOltdLCJpc2FkbSI6dHJ1ZX0.G4MZzNb2G7DF48y2kCFNNYTXyRluQAX2zpSbNmBNk6SVJqVMyUbsCbTlwYV-o7frCDR2x28EYqJBxRiTVzaMiV-QfTz5k8cI64zhWlbV1WYycnIPfgdJkY727ZiMLKqB3D8YTuunZqZjNdRdkUnzY3g2K8M2YJc32Yfn6qRrJ51_bFM2UhcEQVg7P0V6qGw273lsiXJjjFDHDHoFW6FENUMv7Q_4gfcLRJz2BwKe2Ug9jAoO5WEuAPRlBvQeLpGk7S6_-fbE8sdygChjYQVcwBxER7mqlud2VCCf6eFH3XCjqJNH4VhLW6erIyO9PiWzyeSCreujzlpPLXoh6vCGbA",
			},
			want: expected{claims: &someUserClaims, ctxClaim: userCtxClaim, err: nil},
		},
		{
			name: "error token expired",
			args: args{
				header: headerToken,
				token:  "eyJ0eXAiOiJKV1QiLCJhbGciOiJSUzI1NiJ9.eyJyb2xlIjpbInNvbWVSb2xlIl0sImVtYWlsIjoiZXhhbXBsZUBwb21lbG8ubGEiLCJhcmVhIjpbInNvbWVCVSJdLCJ0ZWFtcyI6WyJzb21lVGVhbSJdLCJleHAiOjE2NDM1NzM1MzUsImF1ZCI6InVzZXJjb250ZXh0YXBpIiwic2xhY2siOiJhYmMxMjMiLCJkZWxlZ2F0ZWQiOnRydWUsIm1hbmFnZXIiOm51bGwsImlkIjoiMTIzNGFjYmQiLCJzY29wZXMiOltdLCJpc2FkbSI6dHJ1ZX0.Al638hElxAfdUPoJvFijGSNNKgnMuMUKv2tQNIdv_1QDhKvRHVwmTeK0B1_k1_a0c2DK45l_Bx-LvQEe211kDRZT3T-578XvD4uS4Rhk-G3raVnjMAIyqmLApIvq0aTAWqXx_rCLUIhJO_Y5gUDzsarFetnqvcvfJ8XITJPTX9tP2fWgbzuj-ReNuqOAKDgnQl7jvow-bEkKgKgIHPiU1VHc6PWG0ZPpVcd1HPw744FSopgpkLavawPnD6drFNOcas08-cgCj21hKpiCDfxuQ1eGft0sbXPtL3ADwZI0KfbCpBI5sfov8lJ315bEf0X393eCEgytWj3HS5gZzvB_-w",
			},
			want: expected{err: auth.ErrUnauthorized},
		},
		{
			name: "error token invalid missing email",
			args: args{
				header: headerToken,
				token:  "eyJ0eXAiOiJKV1QiLCJhbGciOiJSUzI1NiJ9.eyJyb2xlIjpbInNvbWVSb2xlIl0sImFyZWEiOlsic29tZUJVIl0sInRlYW1zIjpbInNvbWVUZWFtIl0sImV4cCI6MzcwMDcyMjg3NSwiYXVkIjoidXNlcmNvbnRleHRhcGkiLCJzbGFjayI6ImFiYzEyMyIsImRlbGVnYXRlZCI6dHJ1ZSwibWFuYWdlciI6bnVsbCwiaWQiOiIxMjM0YWNiZCIsInNjb3BlcyI6W10sImlzYWRtIjp0cnVlfQ.aPkOhQcMEoPZLLg6Mjet7_eDCqzJl-YY7qXEgDX-vzgenwtK2UZSgvvYKIzP784SUnF4UhJ0_t8hifMik8xs7Fkdrkq5yYA24TqXgZeZnv0FNzp5ixM5EuqgjRawAWZYfTSlhL0SQjwPM5EXFuBlYWcZJnBYYi_kAFU1JTCxzSypWa5km5LJmXjgxlq6A3GqyGrPq-jnmIj-OeQ8YeAnpNL4Z_PzxyjCmucv8lkMr6WAQFuxFjAo7pik5L__QWsqbhSChyA82-Dh2G3__ZsgRIJ9VJkbdBKTt8G60qYcOuz3LU7agULKn23Zbn2xXkCyXpCkUz1xZ6w4utvmpWuP4g",
			},
			want: expected{err: auth.ErrMissingEmail},
		},
		{
			name: "error token invalid missing areas",
			args: args{
				header: headerToken,
				token:  "eyJ0eXAiOiJKV1QiLCJhbGciOiJSUzI1NiJ9.eyJyb2xlIjpbInNvbWVSb2xlIl0sImVtYWlsIjoiZXhhbXBsZUBwb21lbG8ubGEiLCJ0ZWFtcyI6WyJzb21lVGVhbSJdLCJleHAiOjI1MzI0Mzc1ODYsImF1ZCI6InVzZXJjb250ZXh0YXBpIiwic2xhY2siOiJhYmMxMjMiLCJkZWxlZ2F0ZWQiOnRydWUsIm1hbmFnZXIiOm51bGwsImlkIjoiMTIzNGFjYmQiLCJzY29wZXMiOltdLCJpc2FkbSI6dHJ1ZX0.Dt-lVMR3qNZXXci1lEGgpYHiySV-vhFDGFD2_SW_BNWZdrRCkene0kWJAymOQFyDmL-JfUNHIQXaO-OEo6KRxy7XtuUEwTiHmzqPwwdsB_AUVVciLwnjhuKAUdjPjF7pUaq_zMElw2miMKkOldmK0SKLI3Tq7uCyJAyQBNH9Bp1JQBErttVE4QR6FwMI0nX1Eitqz6-x8ewkEOQHeX8mbZqjjTDmFSciJT36pinXHOxZphG88s_meTta_Byq7U_emEJXBgDwcdRuDWyUr21-kIjfznFCNwfTEihDyX8mUBTXTOqxToi63RFxFlPWrZzl7zmLgJIcL-sYaDVJyIgwuQ",
			},
			want: expected{err: auth.ErrMissingAreas},
		},
		{
			name: "error invalid signature",
			args: args{
				header: headerToken,
				token:  "eyJ0eXAiOiJKV1QiLCJhbGciOiJSUzI1NiJ9.eyJyb2xlIjpbImRldmVsb3BlciJdLCJlbWFpbCI6ImV4YW1wbGVAcG9tZWxvLmxhIiwidGVhbXMiOlsic29tZXRlYW0iXSwiYXVkIjoidXNlcmNvbnRleHRhcGkiLCJzbGFjayI6IjEyMzRhYmNkIiwiZGVsZWdhdGVkIjp0cnVlLCJtYW5hZ2VyIjpudWxsLCJpZCI6InF3ZXJ0eWlkIiwic2NvcGVzIjpbXSwiaXNhZG0iOnRydWV9.URKlCjx2ACi3c8yfBKWsWL6j4WHLyvQG02O74ldgUlXRlqHRrvAS7Z5dAzm7ELIDT1U5XZOhJchuin-2HJ9YUtSqHng5H1xgnXDUV6GTa1DFF0rz5xj9TXmQOLUC26fK5ZpWygwC9-D4YIyJu9qCQ1GywhOJC5ZXZLunsYC113AAWQPHAhl4nEMp735HGvf1HB2AK2m_veJg5Hrvx8QlrOuAYTc-FU1k0Veu0tnpUTyeROp0uVO-Tgl4gnJwxJ5yN7JaKV31XYU4jwZ6yfhnLmjq_F98rIdORsRMsfRxrDlvi-lK5hsm9URir_8jlAMYtZMWqXUOlHUbGW_ZgcQGSw",
			},
			want: expected{err: auth.ErrUnauthorized},
		},
		{
			name: "missing header",
			args: args{header: ""},
			want: expected{claims: nil, err: auth.ErrRequestNotAcceptable},
		},
	}

	for _, tt := range tests {
		err := os.Setenv("CONTEXT_API_PUBLIC_KEY", TestRSAPublicKey)
		assert.NoError(t, err)

		t.Run(tt.name, func(t *testing.T) {
			handler := func(w http.ResponseWriter, r *http.Request) {}

			req := httptest.NewRequest(http.MethodGet, "http://test", nil)
			if tt.args.header != "" {
				req.Header.Add(tt.args.header, tt.args.token)
			}

			res := httptest.NewRecorder()
			handler(res, req)

			// When
			got, err := auth.DecodeToken(req)

			// Then
			if tt.want.err != nil {
				assert.ErrorIs(t, tt.want.err, err)
			}
			if tt.want.claims != nil {
				assert.Equal(t, tt.want.claims.Area, got.Area)
				assert.Equal(t, strings.Join(tt.want.claims.Area, ","), req.Header.Get(auth.BusinessUnits))
			}
			// Then user context claims
			if tt.want.claims != nil && tt.want.ctxClaim == userCtxClaim {
				assert.Equal(t, tt.want.claims.Role[0], got.Role[0])
				assert.Equal(t, tt.want.claims.Role[0], req.Header.Get(auth.Role))
				assert.Equal(t, tt.want.claims.Email, got.Email)
				assert.Equal(t, tt.want.claims.Email, req.Header.Get(auth.Owner))
			}
		})
	}
}

func TestDecodeServiceToken(t *testing.T) {
	type args struct {
		header string
		token  string
	}
	type mockClaims auth.Claims
	type expected struct {
		claims   *mockClaims
		ctxClaim string
		err      error
	}
	someSvcClaims := mockClaims{Area: []string{"someBU"}, ServiceName: "mock-api", Role: []string{"service_mock-api"}}
	headerToken := "X-Auth-Token"
	svcCtxClaim := "service"

	tests := []struct {
		name string
		args args
		want expected
	}{
		{
			name: "success service ctx audience",
			args: args{
				header: headerToken,
				token:  "eyJ0eXAiOiJKV1QiLCJhbGciOiJSUzI1NiJ9.eyJzZXJ2aWNlbmFtZSI6Im1vY2stYXBpIiwiYXJlYSI6WyJzb21lQlUiXSwieF9hcHBfaWQiOiJtb2NrLWFwaToxMi1hYi0zNC1jZC0xYTJiM2M0ZCIsImV4cCI6MzcwMDcyMjg3NSwiYXVkIjoic2VydmljZWNvbnRleHRhcGkifQ.asVIYw0DjVJt6mivgZnKsYKhjaztRnK3Z4xlEfZ6T80Q-a1Yb1iojC3CXQs1jeFJySeTkSkVpAX9MeZtVwCFqcycUtXj7MNN5SmTPNLtPe61VLX1ALEWvhKvKQDckWn80QUQXMmGb7vxgknt78Y2hQZBTeunVhCEgHg8fjzk8uFykKlfT3xukVCoh1tTB3HBMPjhi27-pSU1mD6zoKO1jzayuoBTIm-tr6IJ-d0WltJZnvyt5xRnQ00A2gnsRXwdUX7ctV5nFJ-kwUCMh4sEkF7mII6hr7D_3nh1wMsH5fTjQozjhOh4YBYD5ZTgZkrp3FBF0W6vzBkGRCDwxawwLA",
			},
			want: expected{claims: &someSvcClaims, ctxClaim: svcCtxClaim, err: nil},
		},
	}

	for _, tt := range tests {
		err := os.Setenv("CONTEXT_API_PUBLIC_KEY", TestRSAPublicKey)
		assert.NoError(t, err)

		t.Run(tt.name, func(t *testing.T) {
			handler := func(w http.ResponseWriter, r *http.Request) {}

			req := httptest.NewRequest(http.MethodGet, "http://test", nil)
			if tt.args.header != "" {
				req.Header.Add(tt.args.header, tt.args.token)
			}

			res := httptest.NewRecorder()
			handler(res, req)

			// When
			got, err := auth.DecodeToken(req)

			// Then
			if tt.want.err != nil {
				assert.ErrorIs(t, tt.want.err, err)
			}
			if tt.want.claims != nil {
				assert.Equal(t, tt.want.claims.Area, got.Area)
				assert.Equal(t, strings.Join(tt.want.claims.Area, ","), req.Header.Get(auth.BusinessUnits))
			}
			// Then service context claims
			if tt.want.claims != nil && tt.want.ctxClaim == svcCtxClaim {
				assert.Equal(t, tt.want.claims.ServiceName, got.ServiceName)
				assert.Equal(t, tt.want.claims.ServiceName, req.Header.Get(auth.Owner))
				assert.Equal(t, tt.want.claims.Role[0], req.Header.Get(auth.Role))
			}
		})
	}
}
