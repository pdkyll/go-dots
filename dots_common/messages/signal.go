package messages

import "fmt"

type TargetPortRange struct {
	LowerPort int `json:"lower-port" cbor:"lower-port"`
	UpperPort int `json:"upper-port" cbor:"upper-port"`
}

type MitigationScope struct {
	Scopes []Scope `json:"scope" cbor:"scope"`
}

type Scope struct {
	// Identifier for the mitigation request
	MitigationId int `json:"mitigation-id" cbor:"mitigation-id"`
	// IP address
	TargetIp []string `json:"target-ip" cbor:"target-ip"`
	// prefix
	TargetPrefix []string `json:"target-prefix" cbor:"target-prefix"`
	// lower-port upper-port
	TargetPortRange []TargetPortRange `json:"target-port-range" cbor:"target-port-range"`
	// Internet Protocol number
	TargetProtocol []int `json:"target-protocol" cbor:"target-protocol"`
	// FQDN
	FQDN []string `json:"FQDN" cbor:"FQDN"`
	// URI
	URI []string `json:"URI" cbor:"URI"`
	// alias name
	AliasName []string `json:"alias-name" cbor:"alias-name"`
	// lifetime
	Lifetime int `json:"lifetime" cbor:"lifetime"`
}

type MitigationRequest struct {
	MitigationScope MitigationScope `json:"mitigation-scope" cbor:"mitigation-scope"`
}

/*
 * Convert MitigationRequests to strings
 */
func (m *MitigationRequest) String() (result string) {
	result = "\n"
	for key, scope := range m.MitigationScope.Scopes {
		result += fmt.Sprintf("   \"%s[%d]\":\n", "scope", key+1)
		result += fmt.Sprintf("     \"%s\": %d\n", "mitigation-id", scope.MitigationId)
		if scope.TargetIp != nil {
			for k, v := range scope.TargetIp {
				result += fmt.Sprintf("     \"%s[%d]\": %s\n", "target-ip", k+1, v)
			}
		}
		if scope.TargetPrefix != nil {
			for k, v := range scope.TargetPrefix {
				result += fmt.Sprintf("     \"%s[%d]\": %s\n", "target-prefix", k+1, v)
			}
		}
		if scope.TargetPortRange != nil {
			for k, v := range scope.TargetPortRange {
				result += fmt.Sprintf("     \"%s[%d]\":\n", "target-port-range", k+1)
				result += fmt.Sprintf("       \"%s\": %d\n", "lower-port", v.LowerPort)
				result += fmt.Sprintf("       \"%s\": %d\n", "upper-port", v.UpperPort)
			}
		}
		if scope.TargetProtocol != nil {
			for k, v := range scope.TargetProtocol {
				result += fmt.Sprintf("     \"%s[%d]\": %d\n", "target-protocol", k+1, v)
			}
		}
		if scope.FQDN != nil {
			for k, v := range scope.FQDN {
				result += fmt.Sprintf("     \"%s[%d]\": %s\n", "FQDN", k+1, v)
			}
		}
		if scope.URI != nil {
			for k, v := range scope.URI {
				result += fmt.Sprintf("     \"%s[%d]\": %s\n", "URI", k+1, v)
			}
		}
		if scope.AliasName != nil {
			for k, v := range scope.AliasName {
				result += fmt.Sprintf("     \"%s[%d]\": %s\n", "alias-name", k+1, v)
			}
		}
		result += fmt.Sprintf("     \"%s\": %d\n", "lifetime", scope.Lifetime)
	}
	return
}

type SignalConfig struct {
	// Identifier for the DOTS signal channel session configuration data represented as an integer.
	// This identifier MUST be generated by the DOTS client.  This document does not make any assumption about how this
	// identifier is generated. This is a mandatory attribute.
	SessionId int `json:"session-id" cbor:"session-id"`
	// Heartbeat interval to check the DOTS peer health.  This is an optional attribute.
	HeartbeatInterval int `json:"heartbeat-interval" cbor:"heartbeat-interval"`
	// Maximum number of missing heartbeat response allowed. This is an optional attribute.
	MissingHbAllowed int `json:"missing-hb-allowed" cbor:"missing-hb-allowed"`
	// Maximum number of retransmissions for a message (referred to as MAX_RETRANSMIT parameter in CoAP).
	// This is an optional attribute.
	MaxRetransmit int `json:"max-retransmit" cbor:"max-retransmit"`
	// Timeout value in seconds used to calculate the initial retransmission timeout value (referred to as ACK_TIMEOUT
	// parameter in CoAP). This is an optional attribute.
	AckTimeout int `json:"ack-timeout" cbor:"ack-timeout"`
	// Random factor used to influence the timing of retransmissions (referred to as ACK_RANDOM_FACTOR parameter in
	// CoAP).  This is an optional attribute.
	AckRandomFactor float64 `json:"ack-random-factor" cbor:"ack-random-factor"`
	// If false, mitigation is triggered only if the signal channel is lost. This is an optional attribute.
	TriggerMitigation bool `json:"trigger-mitigation" cbor:"trigger-mitigation"`
}

type HelloRequest struct {
	Message string `json:"message" cbor:"message"`
}

type HelloResponse struct {
	Message string `json:"message" cbor:"message"`
}
