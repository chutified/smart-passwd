package handlers

import (
	"net/http"

	"github.com/chutified/smart-passwd/config"
	"github.com/chutified/smart-passwd/controls"
	"github.com/chutified/smart-passwd/models"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

// PWDhandler handles the backend generation.
type PWDhandler struct {
	pwdCtrl *controls.Controller
}

// NewPWD is the constructor of the PWDhandler.
func NewPWD() *PWDhandler {
	return &PWDhandler{
		pwdCtrl: controls.New(),
	}
}

// Init starts the controller's services.
func (h *PWDhandler) Init(cfg *config.DBConfig) error {
	// init the controller
	err := h.pwdCtrl.Init(cfg)
	if err != nil {
		return errors.Wrap(err, "initializing password controller")
	}

	return nil
}

// Close stops all connections.
func (h *PWDhandler) Close() error {
	// stop the controller
	err := h.pwdCtrl.Stop()
	if err != nil {
		return errors.Wrap(err, "stopping password controller")
	}

	return nil
}

// PasswordGen handles the password generation.
func (h *PWDhandler) PasswordGen(c *gin.Context) {
	// bind JSON
	var preq models.PasswordReq
	if err := c.ShouldBindJSON(&preq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})

		return
	}

	// generate the passwd
	resp, err := h.pwdCtrl.Generate(&preq)
	if errors.Is(err, controls.ErrInvalidLen) {
		// both length and helper are missing
		c.JSON(http.StatusBadRequest,
			gin.H{"error": err.Error()})

		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})

		return
	}

	// response
	c.JSON(200, resp)
}
