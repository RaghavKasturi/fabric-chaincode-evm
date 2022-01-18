/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package eventmanager

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/hyperledger/fabric-chaincode-evm/evmcc/event"
	"github.com/hyperledger/fabric/common/flogging"
	"github.com/hyperledger/fabric-chaincode-go/shim"

	"github.com/hyperledger/burrow/execution/errors"
	"github.com/hyperledger/burrow/execution/evm"
	"github.com/hyperledger/burrow/execution/exec"
)

var logger = flogging.MustGetLogger("eventmanager")

type EventManager struct {
	Stub       shim.ChaincodeStubInterface
	EventCache []event.Event
	Context    map[string]string
}

var _ evm.EventSink = &EventManager{}

type EventSink interface {
	Call(call *exec.CallEvent, exception *errors.Exception) error
	Log(log *exec.LogEvent) error
}

// Flush will marshal all collected events from the transaction
// and set as a singular Fabric event
//
// eventName is for fabric, typically the evm 8byte function hash
func (evmgr *EventManager) Flush(eventName string) error {
	if len(evmgr.EventCache) == 0 {
		logger.Infof("No events to flush")
		return nil
	}
	payload, err := json.Marshal(evmgr.EventCache)
	if err != nil {
		logger.Errorf("Failed to marshal events: %v", err)
		return fmt.Errorf("Failed to marshal event messages: %s", err)
	}
	logger.Infof("Flushing events for %v", eventName)
	return evmgr.Stub.SetEvent(eventName, payload)
}

// Call for right now is a noop
// need to figure out what it means (burrow or evm)
func (evmgr *EventManager) Call(call *exec.CallEvent, exception *errors.Exception) error {
	return nil
}

// Log will take the given log message convert to a event type and
// append to the event manager's EventCache
func (evmgr *EventManager) Log(log *exec.LogEvent) error {
	e := event.Event{
		Address: strings.ToLower(log.Address.String()),
		Data:    strings.ToLower(log.Data.String()),
		Context: evmgr.Context,
	}
	logger.Infof("Logging event %v", e.Address)

	var topics []string
	for _, topic := range log.Topics {
		t, err := topic.MarshalText()
		if err != nil {
			return fmt.Errorf("Failed to Marshal Topic: %s", err)
		}
		topics = append(topics, strings.ToLower(string(t)))
	}
	e.Topics = topics

	evmgr.EventCache = append(evmgr.EventCache, e)

	return nil
}
