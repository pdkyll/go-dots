package main

import (
    "bytes"
    "errors"
    "net"
    "reflect"

    log "github.com/sirupsen/logrus"
    "github.com/gonuts/cbor"
    "github.com/ugorji/go/codec"

    "github.com/nttdots/go-dots/dots_common/messages"
    "github.com/nttdots/go-dots/dots_server/controllers"
    "github.com/nttdots/go-dots/dots_server/models"
    "github.com/nttdots/go-dots/libcoap"
)

type DotsServiceMethod func(request interface{}, customer *models.Customer) (controllers.Response, error)

func unmarshalCbor(pdu *libcoap.Pdu, typ reflect.Type) (interface{}, error) {
    if len(pdu.Data) == 0 {
        return nil, nil
    }

    m := reflect.New(typ).Interface()
    reader := bytes.NewReader(pdu.Data)

    handle := new(codec.CborHandle)
    d := codec.NewDecoder(reader, handle)
    err := d.Decode(m)

    if err != nil {
        return nil, err
    }
    return m, nil
}

func marshalCbor(msg interface{}) ([]byte, error) {
    writer := bytes.NewBuffer(nil)
    e := cbor.NewEncoder(writer)

    err := e.Encode(msg)
    if err != nil {
        return nil, err
    }
    return writer.Bytes(), nil
}

func addHandler(ctx *libcoap.Context, code messages.Code, controller controllers.ControllerInterface) {
    msg := messages.MessageTypes[code]

    path := "/" + msg.Path
    resource := libcoap.ResourceInit(&path, 0)
    var toMethodHandler = func(method DotsServiceMethod) libcoap.MethodHandler {
        return func(context  *libcoap.Context,
                    resource *libcoap.Resource,
                    session  *libcoap.Session,
                    request  *libcoap.Pdu,
                    token    *[]byte,
                    query    *string,
                    response *libcoap.Pdu) {

            log.WithField("MessageID", request.MessageID).Info("Incoming Request")


            response.MessageID = request.MessageID
            response.Token     = request.Token

            cn, err := session.DtlsGetPeerCommonName()
            if err != nil {
                log.WithError(err).Warn("DtlsGetPeercCommonName() failed")
                response.Code = libcoap.ResponseForbidden
                return
            }

            log.Infof("CommonName is %v", cn)

            customer, err := models.GetCustomerByCommonName(cn)
            if err != nil || customer.Id == 0 {
                log.WithError(err).Warn("Customer not found.")
                response.Code = libcoap.ResponseForbidden
                return
            }

            req, err := unmarshalCbor(request, msg.Type)
            if err != nil {
                log.WithError(err).Error("unmarshalCbor failed.")
                response.Code = libcoap.ResponseInternalServerError
                return
            }

            res, err := method(req, customer)
            if err != nil {
                log.WithError(err).Error("controller returned error")
                response.Code = libcoap.ResponseInternalServerError
                return
            }

            payload, err := marshalCbor(res.Body)
            if err != nil {
                log.WithError(err).Error("marshalCbor failed.")
                response.Code = libcoap.ResponseInternalServerError
                return
            }

            response.Code = libcoap.Code(res.Code)
            response.Data = payload

            return
        }
    }

    resource.RegisterHandler(libcoap.RequestGet,    toMethodHandler(controller.Get))
    resource.RegisterHandler(libcoap.RequestPut,    toMethodHandler(controller.Put))
    resource.RegisterHandler(libcoap.RequestPost,   toMethodHandler(controller.Post))
    resource.RegisterHandler(libcoap.RequestDelete, toMethodHandler(controller.Delete))
    ctx.AddResource(resource)
}

func listen(address string, port uint16, dtlsParam *libcoap.DtlsParam) (_ *libcoap.Context, err error) {
    log.Debugf("listen.go, listen -in. address=%+v, port=%+v", address, port)
    ip := net.ParseIP(address)
    if ip == nil {
        err = errors.New("net.ParseIP() -> nil")
        return
    }

    addr, err := libcoap.AddressOf(ip, port)
    if err != nil {
        return
    }
    log.Debugf("addr=%+v", addr)

    ctx := libcoap.NewContextDtls(nil, dtlsParam)
    if ctx == nil {
        err = errors.New("libcoap.NewContextDtls() -> nil")
        return
    }

    ctx.NewEndpoint(addr, libcoap.ProtoDtls)
    return ctx, nil
}


func listenData(address string, port uint16, dtlsParam *libcoap.DtlsParam) (_ *libcoap.Context, err error) {
    ctx, err := listen(address, port, dtlsParam)
    if err != nil {
        return
    }

    addHandler(ctx, messages.HELLO,                  &controllers.Hello{})
    addHandler(ctx, messages.CREATE_IDENTIFIERS,     &controllers.CreateIdentifiers{})
    addHandler(ctx, messages.INSTALL_FILTERING_RULE, &controllers.InstallFilteringRule{})

    return ctx, nil
}

func listenSignal(address string, port uint16, dtlsParam *libcoap.DtlsParam) (_ *libcoap.Context, err error) {
    ctx, err := listen(address, port, dtlsParam)
    if err != nil {
        return
    }

    addHandler(ctx, messages.HELLO,                 &controllers.Hello{})
    addHandler(ctx, messages.MITIGATION_REQUEST,    &controllers.MitigationRequest{})
    addHandler(ctx, messages.SESSION_CONFIGURATION, &controllers.SessionConfiguration{})

    return ctx, nil
}