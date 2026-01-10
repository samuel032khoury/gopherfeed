{{define "subject"}}Finish Registration with Gopherfeed{{end}}

{{define "body"}}
<!doctype html>
<html>
    <head>
        <meta name="viewport" content="width=device-width" /> 
        <meta http-equiv="Content-Type" content="text/html; charset=UTF-8" />
        <title>Finish Registration</title>
    </head>
    <body>
        <p>Hi {{.Username}},</p>
        <p>Thanks for signing up for Gopherfeed. We're excited to have you on board!</p>
        <p>Before you can start using your account, please complete your registration by clicking the link below to confirm your email address:</p>
        <p><a href="{{.ActivationURL}}">Complete Registration</a></p>
        <p>Or copy and paste the following URL into your web browser:</p>
        <p>{{.ActivationURL}}</p>
        <p>If you did not sign up for a Gopherfeed account, please ignore this email.</p>

        <p>Cheers,</p>
        <p>The Gopherfeed Team</p>
    </body>
</html>
{{end}}