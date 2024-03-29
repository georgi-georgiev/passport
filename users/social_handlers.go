package users

import (
	"net/http"

	"github.com/georgi-georgiev/blunder"
	"github.com/georgi-georgiev/passport/responses"
	"github.com/gin-gonic/gin"
	fb "github.com/huandu/facebook/v2"
)

type FacebookUser struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

// FacebookCallbackHandler godoc
// @Summary Facebook Callback
// @Description Facebook callback
// @Tags identity
// @Accept  json
// @Produce  json
// @Param data body CreateUserPayload true "data"
// @Success 201 {object} CreateUserResponse
// @Router /facebook/callback [post]
func (h *UserHandlers) FacebookCallback(c *gin.Context) {
	accessToken := c.PostForm("accessToken")

	facebookUser, err := verifyFacebookToken(accessToken)
	if err != nil {
		h.blunder.GinAdd(c, err)
		return
	}

	existingUser, err := h.userService.GetUserByEmail(c.Request.Context(), facebookUser.Email)
	if err != nil {
		c.JSON(http.StatusBadRequest, blunder.BadRequest())
		return
	}

	if existingUser != nil {

		token, exp, err := h.userService.IssueAccessToken(existingUser)
		if err != nil {
			h.blunder.GinAdd(c, err)
			return
		}

		refreshToken, err := h.userService.IssueRefreshToken(existingUser)
		if err != nil {
			h.blunder.GinAdd(c, err)
			return
		}

		c.JSON(http.StatusOK, responses.CreateUserResponse{ID: existingUser.ID.Hex(), TokenType: "Bearer", AccessToken: token, RefreshToken: refreshToken, ExpiresIn: exp})
		return
	}

	user, err := h.userService.CreateUser(c.Request.Context(), "", facebookUser.Email, "", "basic", false, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, blunder.InternalServerError())
		return
	}

	token, exp, err := h.userService.IssueAccessToken(user)
	if err != nil {
		h.blunder.GinAdd(c, err)
		return
	}

	refreshToken, err := h.userService.IssueRefreshToken(user)
	if err != nil {
		h.blunder.GinAdd(c, err)
		return
	}

	c.JSON(http.StatusCreated, responses.CreateUserResponse{ID: existingUser.ID.Hex(), TokenType: "Bearer", AccessToken: token, RefreshToken: refreshToken, ExpiresIn: exp})
}

func verifyFacebookToken(accessToken string) (*FacebookUser, error) {
	res, err := fb.Get("/me", fb.Params{
		"fields":       "id,name,email",
		"access_token": accessToken,
	})

	if err != nil {
		return nil, err
	}

	user := &FacebookUser{}
	err = res.Decode(user)
	if err != nil {
		return nil, err
	}

	return user, nil
}
