{{define "auth/activate_email"}}
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>邮箱验证码(Email verification code)</title>
    <style>
        .code {
            color: blue;
            text-decoration: underline;
            display: inline-block;
        }
    </style>
</head>
<body>
    <span>验证码：<h3 class="code">{{.Code}}</h3></span>
    <p>此验证码用于注册验证邮箱使用，有效期为 {{.ActiveCodeLives}} 若非本人操作，请忽略该邮件!</p>
</body>
</html>
{{end}}