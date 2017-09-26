package models

import (
    "time"
    "github.com/nttdots/go-dots/dots_server/db_models"
)

type AcctV5 struct {
    AgentId         int
    ClassId         string
    MacSrc          string
    MacDst          string
    Vlan            int
    IpSrc           string
    IpDst           string
    SrcPort         int
    DstPort         int
    IpProto         string
    Tos             int
    Packets         int
    Bytes           int64
    Flows           int
    StampInserted   time.Time
    StampUpdated    time.Time
}

func NewAcctV5() (s *AcctV5) {
    s = &AcctV5{
        0,
        "",
        "",
        "",
        0,
        "",
        "",
        0,
        0,
        "",
        0,
        0,
        0,
        0,
        time.Unix(0, 0),
        time.Unix(0, 0),
    }
    return
}

func CreateAcctV5Model(dbAcctV5 *db_models.AcctV5) (*AcctV5) {
    acctV5 := &AcctV5{}

    acctV5.AgentId = dbAcctV5.AgentId
    acctV5.ClassId = dbAcctV5.ClassId
    acctV5.MacSrc = dbAcctV5.MacSrc
    acctV5.MacDst = dbAcctV5.MacDst
    acctV5.Vlan = dbAcctV5.Vlan
    acctV5.IpSrc = dbAcctV5.IpSrc
    acctV5.IpDst = dbAcctV5.IpDst
    acctV5.SrcPort = dbAcctV5.SrcPort
    acctV5.DstPort = dbAcctV5.DstPort
    acctV5.IpProto = dbAcctV5.IpProto
    acctV5.Tos = dbAcctV5.Tos
    acctV5.Packets = dbAcctV5.Packets
    acctV5.Bytes = dbAcctV5.Bytes
    acctV5.Flows = dbAcctV5.Flows
    acctV5.StampInserted = dbAcctV5.StampInserted
    acctV5.StampUpdated = dbAcctV5.StampUpdated

    return acctV5
}

func CreateAcctV5DbModel(acctV5 *AcctV5) (*db_models.AcctV5) {
    dbAcctV5 := &db_models.AcctV5{}

    dbAcctV5.AgentId = acctV5.AgentId
    dbAcctV5.ClassId = acctV5.ClassId
    dbAcctV5.MacSrc = acctV5.MacSrc
    dbAcctV5.MacDst = acctV5.MacDst
    dbAcctV5.Vlan = acctV5.Vlan
    dbAcctV5.IpSrc = acctV5.IpSrc
    dbAcctV5.IpDst = acctV5.IpDst
    dbAcctV5.SrcPort = acctV5.SrcPort
    dbAcctV5.DstPort = acctV5.DstPort
    dbAcctV5.IpProto = acctV5.IpProto
    dbAcctV5.Tos = acctV5.Tos
    dbAcctV5.Packets = acctV5.Packets
    dbAcctV5.Bytes = acctV5.Bytes
    dbAcctV5.Flows = acctV5.Flows
    dbAcctV5.StampInserted = acctV5.StampInserted
    dbAcctV5.StampUpdated = acctV5.StampUpdated

    return dbAcctV5
}