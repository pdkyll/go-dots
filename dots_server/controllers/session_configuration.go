package controllers

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	common "github.com/nttdots/go-dots/dots_common"
	"github.com/nttdots/go-dots/dots_common/messages"
	"github.com/nttdots/go-dots/dots_server/models"
	dots_config "github.com/nttdots/go-dots/dots_server/config"
)

/*
 * Controller for the session_configuration API.
 */
type SessionConfiguration struct {
	Controller
}

func (m *SessionConfiguration) Get(request interface{}, customer *models.Customer) (res Response, err error) {

	signalSessionConfiguration, err := models.GetCurrentSignalSessionConfiguration(customer.Id)
	if err != nil {
		res = Response{
			Type: common.NonConfirmable,
			Code: common.BadRequest,
			Body: nil,
		}
		return
	}

	// TODO: check found or not

	config := dots_config.GetServerSystemConfig().SignalConfigurationParameter

	resp := messages.ConfigurationResponse{}
	resp.HeartbeatInterval.SetMinMax(config.HeartbeatInterval)
	resp.MissingHbAllowed.SetMinMax(config.MissingHbAllowed)
	resp.MaxRetransmit.SetMinMax(config.MaxRetransmit)
	resp.AckTimeout.SetMinMax(config.AckTimeout)
	resp.AckRandomFactor.SetMinMax(config.AckRandomFactor)

	resp.HeartbeatInterval.CurrentValue = signalSessionConfiguration.HeartbeatInterval
	resp.MissingHbAllowed.CurrentValue  = signalSessionConfiguration.MissingHbAllowed
	resp.MaxRetransmit.CurrentValue     = signalSessionConfiguration.MaxRetransmit
	resp.AckTimeout.CurrentValue        = signalSessionConfiguration.AckTimeout
	resp.AckRandomFactor.CurrentValue   = signalSessionConfiguration.AckRandomFactor
	resp.TriggerMitigation.CurrentValue = signalSessionConfiguration.TriggerMitigation

	res = Response{
			Type: common.NonConfirmable,
			Code: common.Content,
			Body: resp,
	}

	return
}

/*
 * Handles session_configuration PUT requests and start the mitigation.
 *  1. Validate the received session configuration requests.
 *  2. return the validation results.
 *
 * parameter:
 *  request request message
 *  customer request source Customer
 * return:
 *  res response message
 *  err error
 */
func (m *SessionConfiguration) Put(request interface{}, customer *models.Customer) (res Response, err error) {

	if request == nil {
		res = Response{
			Type: common.NonConfirmable,
			Code: common.BadRequest,
			Body: nil,
		}
		return
	}

	payload := &request.(*messages.SignalConfigRequest).SignalConfig
	sessionConfigurationPayloadDisplay(payload)

	// validate
	signalSessionConfiguration := models.NewSignalSessionConfiguration(
		payload.SessionId,
		payload.HeartbeatInterval,
		payload.MissingHbAllowed,
		payload.MaxRetransmit,
		payload.AckTimeout,
		payload.AckRandomFactor,
		payload.TriggerMitigation,
	)
	v := models.SignalConfigurationValidator{}
	validateResult := v.Validate(signalSessionConfiguration, *customer)
	if !validateResult {
		goto ResponseNG
	} else {
		// Register SignalConfigurationParameter
		_, err = models.CreateSignalSessionConfiguration(*signalSessionConfiguration, *customer)
		if err != nil {
			goto ResponseNG
		}

		goto ResponseOK
	}

ResponseNG:
// on validation error
	res = Response{
		Type: common.NonConfirmable,
		Code: common.BadRequest,
		Body: nil,
	}
	return
ResponseOK:
// on validation success
	res = Response{
		Type: common.NonConfirmable,
		Code: common.Created,
		Body: nil,
	}
	return
}

func (m *SessionConfiguration) Delete(request interface{}, customer *models.Customer) (res Response, err error) {
	err = models.DeleteSignalSessionConfigurationByCustomerId(customer.Id)
	if err != nil {
		res = Response{
			Type: common.NonConfirmable,
			Code: common.InternalServerError,
			Body: nil,
		}
		return
	}

	res = Response{
		Type: common.NonConfirmable,
		Code: common.Deleted,
		Body: nil,
	}
	return
}


/*
 * Parse the request body and display the contents of the messages to stdout.
*/
func sessionConfigurationPayloadDisplay(data *messages.SignalConfig) {

	var result string = "\n"
	result += fmt.Sprintf("   \"%s\": %d\n", "session-id", data.SessionId)
	result += fmt.Sprintf("   \"%s\": %d\n", "heartbeat-interval", data.HeartbeatInterval)
	result += fmt.Sprintf("   \"%s\": %d\n", "missing-hb-allowed", data.MissingHbAllowed)
	result += fmt.Sprintf("   \"%s\": %d\n", "max-retransmit", data.MaxRetransmit)
	result += fmt.Sprintf("   \"%s\": %d\n", "ack-timeout", data.AckTimeout)
	result += fmt.Sprintf("   \"%s\": %f\n", "ack-random-factor", data.AckRandomFactor)
	result += fmt.Sprintf("   \"%s\": %f\n", "trigger-mitigation", data.TriggerMitigation)
	log.Infoln(result)
}
