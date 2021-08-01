package todus

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_parse_token(t *testing.T) {
	token_str := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2Mjc4Njk1NzksInVzZXJuYW1lIjoiNTM1NjA5MTE5MCIsInZlcnNpb24iOiIyMTgyMyJ9.5mer_Rf-WCIJL2OwML2p8VxjuaXPiOn0-B-bNb6_OgY"
	token, phone, err := parse_token(token_str)
	fmt.Println(phone)
	require.NotEmpty(t, token)
	require.NotEmpty(t, phone)
	require.Equal(t, phone, "5356091190")
	require.Equal(t, token, token_str)
	require.NoError(t, err)
}

func Test_steal_token(t *testing.T) {
	token, phone, err := steal_token()
	fmt.Println(phone)
	fmt.Println(token)
	require.NotEmpty(t, token)
	require.NotEmpty(t, phone)
	require.NoError(t, err)
}

func Test_Sign_url(t *testing.T) {
	up, down, token, err := sign_url(1050)
	fmt.Println(up)
	fmt.Println(down)
	require.NotEmpty(t, up)
	require.NotEmpty(t, down)
	require.NotEmpty(t, token)
	require.NoError(t, err)
}

func Test_Upload(t *testing.T) {
	down, err := Upload_File("/media/hoax/h0ax/test.mp3")
	fmt.Println(down)
	require.NotEmpty(t, down)
	require.NoError(t, err)
}
