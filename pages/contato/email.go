package contato

import (
	"html"
	"strings"
)

func buildEmailHTML(nome, telefone, email, modalidade, mensagem string) string {
	r := strings.NewReplacer(
		"{{nome}}", html.EscapeString(nome),
		"{{telefone}}", html.EscapeString(telefone),
		"{{email}}", html.EscapeString(email),
		"{{modalidade}}", html.EscapeString(modalidade),
		"{{mensagem}}", strings.ReplaceAll(html.EscapeString(mensagem), "\n", "<br>"),
	)
	return r.Replace(emailTemplate)
}

const emailTemplate = `<!DOCTYPE html>
<html lang="pt">
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width,initial-scale=1.0">
  <title>Novo contacto — Arena VIP</title>
</head>
<body style="margin:0;padding:0;background-color:#e5e5e5;-webkit-text-size-adjust:100%;">
<table role="presentation" cellpadding="0" cellspacing="0" border="0" width="100%" style="background-color:#e5e5e5;">
  <tr>
    <td align="center" valign="top" style="padding:48px 16px;">
      <table role="presentation" cellpadding="0" cellspacing="0" border="0" width="600" style="max-width:600px;width:100%;">

        <!-- TOP GOLD STRIP -->
        <tr><td height="5" style="background-color:#c9a84c;font-size:0;line-height:0;">&nbsp;</td></tr>

        <!-- HEADER -->
        <tr>
          <td style="background-color:#080808;padding:52px 48px 44px;text-align:center;">
            <p style="margin:0 0 22px 0;font-family:Arial,Helvetica,sans-serif;font-size:9px;letter-spacing:7px;color:#c9a84c;text-transform:uppercase;font-weight:bold;">Academia de Artes Marciais</p>
            <!-- Subtle separator above title -->
            <table role="presentation" cellpadding="0" cellspacing="0" border="0" width="100%" style="margin-bottom:22px;">
              <tr>
                <td style="border-top:1px solid #222;">&nbsp;</td>
                <td style="text-align:center;color:#444;font-size:13px;padding:0 14px;white-space:nowrap;font-family:Arial;">◆</td>
                <td style="border-top:1px solid #222;">&nbsp;</td>
              </tr>
            </table>
            <!-- Title -->
            <h1 style="margin:0;font-family:'Arial Black',Arial,Helvetica,sans-serif;font-size:64px;font-weight:900;color:#ffffff;letter-spacing:12px;text-transform:uppercase;line-height:1;">ARENA VIP</h1>
            <!-- Gold separator below title -->
            <table role="presentation" cellpadding="0" cellspacing="0" border="0" width="100%" style="margin-top:28px;">
              <tr>
                <td style="border-top:1px solid #c9a84c;">&nbsp;</td>
                <td style="text-align:center;color:#c9a84c;font-size:15px;padding:0 14px;white-space:nowrap;font-family:Arial;">◆</td>
                <td style="border-top:1px solid #c9a84c;">&nbsp;</td>
              </tr>
            </table>
          </td>
        </tr>

        <!-- RED ALERT BANNER -->
        <tr>
          <td style="background-color:#be1228;padding:15px 48px;text-align:center;">
            <p style="margin:0;font-family:Arial,Helvetica,sans-serif;font-size:10px;letter-spacing:5px;color:#ffffff;text-transform:uppercase;font-weight:bold;">&#9679;&nbsp;&nbsp;Novo pedido de aula experimental</p>
          </td>
        </tr>

        <!-- BODY -->
        <tr>
          <td style="background-color:#0d0d0d;padding:44px 48px 52px;">

            <!-- Intro -->
            <p style="margin:0 0 36px 0;font-family:Arial,Helvetica,sans-serif;font-size:14px;color:#6b7280;line-height:1.75;text-align:center;">Recebeste uma nova inscrição no formulário de contacto.<br>Responde em até <strong style="color:#c9a84c;">24 horas</strong> para confirmar a aula experimental.</p>

            <!-- CONTACT CARD -->
            <table role="presentation" cellpadding="0" cellspacing="0" border="0" width="100%">
              <tr>
                <td style="border-top:2px solid #c9a84c;background-color:#141414;">

                  <!-- Card header -->
                  <table role="presentation" cellpadding="0" cellspacing="0" border="0" width="100%">
                    <tr>
                      <td style="padding:18px 28px;border-bottom:1px solid #1e1e1e;">
                        <p style="margin:0;font-family:Arial,Helvetica,sans-serif;font-size:8px;letter-spacing:4px;color:#c9a84c;text-transform:uppercase;font-weight:bold;">&#9670;&nbsp; Informações de contacto</p>
                      </td>
                    </tr>
                  </table>

                  <!-- Fields -->
                  <table role="presentation" cellpadding="0" cellspacing="0" border="0" width="100%">

                    <!-- Nome -->
                    <tr>
                      <td width="120" valign="middle" style="padding:20px 12px 20px 28px;border-bottom:1px solid #1e1e1e;">
                        <p style="margin:0;font-family:Arial,Helvetica,sans-serif;font-size:8px;letter-spacing:2px;color:#4b5563;text-transform:uppercase;font-weight:bold;">Nome</p>
                      </td>
                      <td valign="middle" style="padding:20px 28px 20px 0;border-bottom:1px solid #1e1e1e;">
                        <p style="margin:0;font-family:'Arial Black',Arial,Helvetica,sans-serif;font-size:17px;color:#ffffff;font-weight:900;letter-spacing:1px;">{{nome}}</p>
                      </td>
                    </tr>

                    <!-- Telefone -->
                    <tr>
                      <td width="120" valign="middle" style="padding:16px 12px 16px 28px;border-bottom:1px solid #1e1e1e;">
                        <p style="margin:0;font-family:Arial,Helvetica,sans-serif;font-size:8px;letter-spacing:2px;color:#4b5563;text-transform:uppercase;font-weight:bold;">Telefone</p>
                      </td>
                      <td valign="middle" style="padding:16px 28px 16px 0;border-bottom:1px solid #1e1e1e;">
                        <p style="margin:0;font-family:Arial,Helvetica,sans-serif;font-size:15px;color:#d1d5db;">{{telefone}}</p>
                      </td>
                    </tr>

                    <!-- Email -->
                    <tr>
                      <td width="120" valign="middle" style="padding:16px 12px 16px 28px;border-bottom:1px solid #1e1e1e;">
                        <p style="margin:0;font-family:Arial,Helvetica,sans-serif;font-size:8px;letter-spacing:2px;color:#4b5563;text-transform:uppercase;font-weight:bold;">Email</p>
                      </td>
                      <td valign="middle" style="padding:16px 28px 16px 0;border-bottom:1px solid #1e1e1e;">
                        <a href="mailto:{{email}}" style="font-family:Arial,Helvetica,sans-serif;font-size:15px;color:#c9a84c;text-decoration:none;font-weight:bold;">{{email}}</a>
                      </td>
                    </tr>

                    <!-- Modalidade -->
                    <tr>
                      <td width="120" valign="middle" style="padding:16px 12px 16px 28px;">
                        <p style="margin:0;font-family:Arial,Helvetica,sans-serif;font-size:8px;letter-spacing:2px;color:#4b5563;text-transform:uppercase;font-weight:bold;">Modalidade</p>
                      </td>
                      <td valign="middle" style="padding:16px 28px 16px 0;">
                        <table role="presentation" cellpadding="0" cellspacing="0" border="0">
                          <tr>
                            <td style="background-color:#be1228;padding:7px 20px;">
                              <p style="margin:0;font-family:Arial,Helvetica,sans-serif;font-size:10px;letter-spacing:3px;color:#ffffff;text-transform:uppercase;font-weight:bold;">{{modalidade}}</p>
                            </td>
                          </tr>
                        </table>
                      </td>
                    </tr>

                  </table>
                </td>
              </tr>
            </table>

            <!-- MENSAGEM -->
            <table role="presentation" cellpadding="0" cellspacing="0" border="0" width="100%" style="margin-top:28px;">
              <tr>
                <td>
                  <p style="margin:0 0 14px 0;font-family:Arial,Helvetica,sans-serif;font-size:8px;letter-spacing:4px;color:#c9a84c;text-transform:uppercase;font-weight:bold;">&#9670;&nbsp; Mensagem</p>
                  <table role="presentation" cellpadding="0" cellspacing="0" border="0" width="100%">
                    <tr>
                      <td style="background-color:#141414;border-left:3px solid #c9a84c;padding:22px 26px;">
                        <p style="margin:0;font-family:Georgia,'Times New Roman',serif;font-size:15px;color:#9ca3af;line-height:1.85;font-style:italic;">{{mensagem}}</p>
                      </td>
                    </tr>
                  </table>
                </td>
              </tr>
            </table>

            <!-- Spacer -->
            <table role="presentation" cellpadding="0" cellspacing="0" border="0" width="100%" style="margin-top:36px;">
              <tr>
                <!-- WhatsApp CTA -->
                <td align="center">
                  <a href="https://wa.me/351934551015" style="display:inline-block;background-color:#128c7e;padding:18px 44px;font-family:Arial,Helvetica,sans-serif;font-size:10px;letter-spacing:4px;color:#ffffff;text-transform:uppercase;font-weight:bold;text-decoration:none;">&#9658;&nbsp; Responder via WhatsApp</a>
                </td>
              </tr>
            </table>

          </td>
        </tr>

        <!-- GOLD RULE -->
        <tr><td height="1" style="background-color:#c9a84c;font-size:0;line-height:0;">&nbsp;</td></tr>

        <!-- FOOTER -->
        <tr>
          <td style="background-color:#080808;padding:32px 48px;text-align:center;">
            <p style="margin:0 0 10px 0;font-family:Arial,Helvetica,sans-serif;font-size:18px;color:#c9a84c;letter-spacing:2px;">◆</p>
            <p style="margin:0 0 8px 0;font-family:'Arial Black',Arial,Helvetica,sans-serif;font-size:13px;letter-spacing:5px;color:#ffffff;text-transform:uppercase;font-weight:bold;">ARENA VIP</p>
            <p style="margin:0 0 14px 0;font-family:Arial,Helvetica,sans-serif;font-size:11px;color:#374151;">Email gerado automaticamente pelo formulário de contacto do website.</p>
            <a href="https://arenavip.pt" style="font-family:Arial,Helvetica,sans-serif;font-size:9px;letter-spacing:3px;color:#c9a84c;text-transform:uppercase;font-weight:bold;text-decoration:none;">arenavip.pt</a>
          </td>
        </tr>

        <!-- BOTTOM GOLD STRIP -->
        <tr><td height="5" style="background-color:#c9a84c;font-size:0;line-height:0;">&nbsp;</td></tr>

      </table>
    </td>
  </tr>
</table>
</body>
</html>`
