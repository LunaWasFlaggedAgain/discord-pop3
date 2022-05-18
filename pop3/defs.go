package pop3

import (
	"net"
	"sync"
)

type Headers struct {
	AcceptLanguage                  string `pop3:"Accept-Language"`
	AlternateRecipient              string `pop3:"Alternate-Recipient"`
	ARCAuthenticationResults        string `pop3:"ARC-Authentication-Results"`
	ARCMessageSignature             string `pop3:"ARC-Message-Signature"`
	ARCSeal                         string `pop3:"ARC-Seal"`
	ArchivedAt                      string `pop3:"Archived-At"`
	AuthenticationResults           string `pop3:"Authentication-Results"`
	AutoSubmitted                   string `pop3:"Auto-Submitted"`
	Autoforwarded                   string `pop3:"Autoforwarded"`
	Autosubmitted                   string `pop3:"Autosubmitted"`
	Bcc                             string `pop3:"Bcc"`
	Cc                              string `pop3:"Cc"`
	Comments                        string `pop3:"Comments"`
	ContentIdentifier               string `pop3:"Content-Identifier"`
	ContentReturn                   string `pop3:"Content-Return"`
	ContentTransferEncoding         string `pop3:"Content-Transfer-Encoding"`
	ContentType                     string `pop3:"Content-Type"`
	Conversion                      string `pop3:"Conversion"`
	ConversionWithLoss              string `pop3:"Conversion-With-Loss"`
	DLExpansionHistory              string `pop3:"DL-Expansion-History"`
	Date                            string `pop3:"Date"`
	DeferredDelivery                string `pop3:"Deferred-Delivery"`
	DeliveryDate                    string `pop3:"Delivery-Date"`
	DiscardedX400IPMSExtensions     string `pop3:"Discarded-X400-IPMS-Extensions"`
	DiscardedX400MTSExtensions      string `pop3:"Discarded-X400-MTS-Extensions"`
	DiscloseRecipients              string `pop3:"Disclose-Recipients"`
	DispositionNotificationOptions  string `pop3:"Disposition-Notification-Options"`
	DispositionNotificationTo       string `pop3:"Disposition-Notification-To"`
	DKIMSignature                   string `pop3:"DKIM-Signature"`
	Encoding                        string `pop3:"Encoding"`
	Encrypted                       string `pop3:"Encrypted"`
	Expires                         string `pop3:"Expires"`
	ExpiryDate                      string `pop3:"Expiry-Date"`
	From                            string `pop3:"From"`
	GenerateDeliveryReport          string `pop3:"Generate-Delivery-Report"`
	Importance                      string `pop3:"Importance"`
	InReplyTo                       string `pop3:"In-Reply-To"`
	IncompleteCopy                  string `pop3:"Incomplete-Copy"`
	Keywords                        string `pop3:"Keywords"`
	Language                        string `pop3:"Language"`
	LatestDeliveryTime              string `pop3:"Latest-Delivery-Time"`
	ListArchive                     string `pop3:"List-Archive"`
	ListHelp                        string `pop3:"List-Help"`
	ListID                          string `pop3:"List-ID"`
	ListOwner                       string `pop3:"List-Owner"`
	ListPost                        string `pop3:"List-Post"`
	ListSubscribe                   string `pop3:"List-Subscribe"`
	Listunsubscribe                 string `pop3:"list-unsubscribe"`
	ListUnsubscribePost             string `pop3:"List-Unsubscribe-Post"`
	MessageContext                  string `pop3:"Message-Context"`
	MessageID                       string `pop3:"Message-ID"`
	MessageType                     string `pop3:"Message-Type"`
	MTPriority                      string `pop3:"MT-Priority"`
	Obsoletes                       string `pop3:"Obsoletes"`
	Organization                    string `pop3:"Organization"`
	OriginalEncodedInformationTypes string `pop3:"Original-Encoded-Information-Types"`
	OriginalFrom                    string `pop3:"Original-From"`
	OriginalMessageID               string `pop3:"Original-Message-ID"`
	Originalrecipient               string `pop3:"original-recipient"`
	OriginatorReturnAddress         string `pop3:"Originator-Return-Address"`
	OriginalSubject                 string `pop3:"Original-Subject"`
	PICSLabel                       string `pop3:"PICS-Label"`
	PreventNonDeliveryReport        string `pop3:"Prevent-NonDelivery-Report"`
	Priority                        string `pop3:"Priority"`
	Received                        string `pop3:"Received"`
	ReceivedSPF                     string `pop3:"Received-SPF"`
	References                      string `pop3:"References"`
	ReplyBy                         string `pop3:"Reply-By"`
	ReplyTo                         string `pop3:"Reply-To"`
	RequireRecipientValidSince      string `pop3:"Require-Recipient-Valid-Since"`
	ResentBcc                       string `pop3:"Resent-Bcc"`
	ResentCc                        string `pop3:"Resent-Cc"`
	ResentDate                      string `pop3:"Resent-Date"`
	ResentFrom                      string `pop3:"Resent-From"`
	ResentMessageID                 string `pop3:"Resent-Message-ID"`
	ResentReplyTo                   string `pop3:"Resent-Reply-To"`
	ResentSender                    string `pop3:"Resent-Sender"`
	ResentTo                        string `pop3:"Resent-To"`
	ReturnPath                      string `pop3:"Return-Path"`
	Sender                          string `pop3:"Sender"`
	Sensitivity                     string `pop3:"Sensitivity"`
	Solicitation                    string `pop3:"Solicitation"`
	Subject                         string `pop3:"Subject"`
	Supersedes                      string `pop3:"Supersedes"`
	TLSReportDomain                 string `pop3:"TLS-Report-Domain"`
	TLSReportSubmitter              string `pop3:"TLS-Report-Submitter"`
	TLSRequired                     string `pop3:"TLS-Required"`
	To                              string `pop3:"To"`
	VBRInfo                         string `pop3:"VBR-Info"`
}

type Email struct {
	Headers    Headers
	RawHeaders string
	Message    string
}

type Pop3Server struct {
	Emails      []*Email
	Listener    net.Listener
	EmailsMutex sync.Mutex
}
