package nbfx

import (
	"io"
)

type Encoder interface {
	Encode(io.Reader) ([]byte, error)
}

type Decoder interface {
	Decode(io.Reader) (string, error)
}

// for MultiByteInt31
const mask_mbi31 = uint32(0x80) //0x80 = 128

const (
	endElement                        byte = 0x01
	comment                           byte = 0x02
	array                             byte = 0x03
	shortAttribute                    byte = 0x04
	attribute                         byte = 0x05
	shortDictionaryAttribute          byte = 0x06
	dictionaryAttribute               byte = 0x07
	shortXmlnsAttribute               byte = 0x08
	xmlnsAttribute                    byte = 0x09
	shortDictionaryXmlnsAttribute     byte = 0x0A
	dictionaryXmlnsAttribute          byte = 0x0B
	prefixDictionaryAttributeA        byte = 0x0C
	prefixAttributeA                  byte = 0x26
	shortElement                      byte = 0x40
	element                           byte = 0x41
	shortDictionaryElement            byte = 0x42
	dictionaryElement                 byte = 0x43
	prefixDictionaryElementA          byte = 0x44
	prefixElementA                    byte = 0x5E
	zeroText                          byte = 0x80
	zeroTextWithEndElement            byte = 0x81
	oneText                           byte = 0x82
	oneTextWithEndElement             byte = 0x83
	falseText                         byte = 0x84
	falseTextWithEndElement           byte = 0x85
	trueText                          byte = 0x86
	trueTextWithEndElement            byte = 0x87
	int8Text                          byte = 0x88
	int8TextWithEndElement            byte = 0x89
	int16Text                         byte = 0x8A
	int16TextWithEndElement           byte = 0x8B
	int32Text                         byte = 0x8C
	int32TextWithEndElement           byte = 0x8D
	int64Text                         byte = 0x8E
	int64TextWithEndElement           byte = 0x8F
	floatText                         byte = 0x90
	floatTextWithEndElement           byte = 0x91
	doubleText                        byte = 0x92
	doubleTextWithEndElement          byte = 0x93
	decimalText                       byte = 0x94
	decimalTextWithEndElement         byte = 0x95
	dateTimeText                      byte = 0x96
	dateTimeTextWithEndElement        byte = 0x97
	chars8Text                        byte = 0x98
	chars8TextWithEndElement          byte = 0x99
	chars16Text                       byte = 0x9A
	chars16TextWithEndElement         byte = 0x9B
	chars32Text                       byte = 0x9C
	chars32TextWithEndElement         byte = 0x9D
	bytes8Text                        byte = 0x9E
	bytes8TextWithEndElement          byte = 0x9F
	bytes16Text                       byte = 0xA0
	bytes16TextWithEndElement         byte = 0xA1
	bytes32Text                       byte = 0xA2
	bytes32TextWithEndElement         byte = 0xA3
	startListText                     byte = 0xA4
	startListTextWithEndElement       byte = 0xA5
	endListText                       byte = 0xA6
	endListTextWithEndElement         byte = 0xA7
	emptyText                         byte = 0xA8
	emptyTextWithEndElement           byte = 0xA9
	dictionaryText                    byte = 0xAA
	dictionaryTextWithEndElement      byte = 0xAB
	uniqueIdText                      byte = 0xAC
	uniqueIdTextWithEndElement        byte = 0xAD
	timeSpanText                      byte = 0xAE
	timeSpanTextWithEndElement        byte = 0xAF
	uuidText                          byte = 0xB0
	uuidTextWithEndElement            byte = 0xB1
	uInt64Text                        byte = 0xB2
	uInt64TextWithEndElement          byte = 0xB3
	boolText                          byte = 0xB4
	boolTextWithEndElement            byte = 0xB5
	unicodeChars8Text                 byte = 0xB6
	unicodeChars8TextWithEndElement   byte = 0xB7
	unicodeChars16Text                byte = 0xB8
	unicodeChars16TextWithEndElement  byte = 0xB9
	unicodeChars32Text                byte = 0xBA
	unicodeChars32TextWithEndElement  byte = 0xBB
	qNameDictionaryText               byte = 0xBC
	qNameDictionaryTextWithEndElement byte = 0xBD
)
