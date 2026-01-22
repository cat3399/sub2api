package service

import (
	"errors"
	"testing"

	"github.com/Wei-Shaw/sub2api/internal/pkg/antigravity"
	"github.com/stretchr/testify/require"
)

func TestApplyAntigravityProjectID_UsesLoadResponseWhenPresent(t *testing.T) {
	tokenInfo := &AntigravityTokenInfo{}
	applyAntigravityProjectID(tokenInfo, "old-project", &antigravity.LoadCodeAssistResponse{
		CloudAICompanionProject: " new-project ",
	}, nil)

	require.False(t, tokenInfo.ProjectIDMissing)
	require.Equal(t, "new-project", tokenInfo.ProjectID)
}

func TestApplyAntigravityProjectID_PreservesExistingWhenLoadFails(t *testing.T) {
	tokenInfo := &AntigravityTokenInfo{}
	applyAntigravityProjectID(tokenInfo, " old-project ", nil, errors.New("loadCodeAssist failed"))

	require.False(t, tokenInfo.ProjectIDMissing)
	require.Equal(t, "old-project", tokenInfo.ProjectID)
}

func TestApplyAntigravityProjectID_MarksMissingWhenNoExistingAndNoLoad(t *testing.T) {
	tokenInfo := &AntigravityTokenInfo{}
	applyAntigravityProjectID(tokenInfo, "", nil, errors.New("loadCodeAssist failed"))

	require.True(t, tokenInfo.ProjectIDMissing)
	require.Equal(t, "", tokenInfo.ProjectID)
}

func TestApplyAntigravityProjectID_MarksMissingWhenLoadEmptyAndNoExisting(t *testing.T) {
	tokenInfo := &AntigravityTokenInfo{}
	applyAntigravityProjectID(tokenInfo, "", &antigravity.LoadCodeAssistResponse{
		CloudAICompanionProject: "  ",
	}, nil)

	require.True(t, tokenInfo.ProjectIDMissing)
	require.Equal(t, "", tokenInfo.ProjectID)
}

func TestApplyAntigravityProjectID_FallsBackToExistingWhenLoadEmpty(t *testing.T) {
	tokenInfo := &AntigravityTokenInfo{}
	applyAntigravityProjectID(tokenInfo, "old-project", &antigravity.LoadCodeAssistResponse{
		CloudAICompanionProject: "",
	}, nil)

	require.False(t, tokenInfo.ProjectIDMissing)
	require.Equal(t, "old-project", tokenInfo.ProjectID)
}

