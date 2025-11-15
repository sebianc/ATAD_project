package services

import (
	"atad_project/models"
	"encoding/xml"
	"fmt"
	"os"
	"strconv"
	"time"
)

// entry point in OFX XML structure
type OFX struct {
	BankMsgs []BankMsg `xml:"BANKMSGSRSV1>STMTTRNRS>STMTRS"`
}

// bank message containing transactions
type BankMsg struct {
	Transactions []OFXTransaction `xml:"BANKTRANLIST>STMTTRN"`
}

// individual transaction structure to match OFX XML
type OFXTransaction struct {
	DatePosted string `xml:"DTPOSTED"`
	Amount     string `xml:"TRNAMT"`
	Name       string `xml:"NAME"`
	Memo       string `xml:"MEMO"`
	Type       string `xml:"TRNTYPE"`
}

// function to import OFX file and convert to transactions from models package
func ImportOFX(filePath string) ([]*models.Transaction, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("cannot read file: %v", err)
	}

	var ofx OFX
	if err := xml.Unmarshal(data, &ofx); err != nil {
		return nil, fmt.Errorf("failed to parse OFX XML: %v", err)
	}

	var transactions []*models.Transaction

	parseDate := func(raw string) string {
		if len(raw) < 8 {
			return raw
		}
		t, err := time.Parse("20060102", raw[:8])
		if err != nil {
			return raw
		}
		return t.Format("2006-01-02")
	}

	for _, msg := range ofx.BankMsgs {
		for _, t := range msg.Transactions {
			amount, err := strconv.ParseFloat(t.Amount, 64)
			if err != nil {
				return nil, fmt.Errorf("invalid amount format: %v", err)
			}

			desc := t.Name
			if t.Memo != "" {
				desc += " - " + t.Memo
			}

			tr := &models.Transaction{
				DATE:        parseDate(t.DatePosted),
				AMOUNT:      amount,
				DESCRIPTION: desc,
				CATEGORY:    t.Type,
			}
			transactions = append(transactions, tr)
		}
	}

	return transactions, nil
}
