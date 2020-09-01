{{define "auth/reset_password"}}
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>修改密码验证码(Reset password verification code)</title>
    <style>
        .code {
            color: blue;
            text-decoration: underline;
            display: inline-block;
        }
    </style>
</head>
<body>
    <span>{{.DisplayName}} 你好！</span>
    <span>验证码：<h3 class="code">{{.Code}}</h3></span>
    <p>此验证码用于修改密码使用，有效期为 {{.ResetPwdCodeLives}} 若非本人操作，请忽略该邮件!</p>
</body>
</html>
{{end}}