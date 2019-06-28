package templates

var Home string = `
<!DOCTYPE html>
<html lang=en>

<head>
  <meta charset=utf-8>
  <meta http-equiv=X-UA-Compatible content="IE=edge">
  <meta name=viewport content="width=device-width,initial-scale=1">
  <link rel=icon href=/favicon.ico> <title>Æternity compiler</title>
	<link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/mini.css/3.0.1/mini-default.min.css">
	<link rel="apple-touch-icon" sizes="180x180" href="https://aeternity.com/images/favicon/apple-touch-icon.png">
	<link rel="icon" type="image/png" href="https://aeternity.com/user/themes/aeon/images/favicon/favicon-32x32.png" sizes="32x32">
	<link rel="icon" type="image/png" href="https://aeternity.com/user/themes/aeon/images/favicon/favicon-16x16.png" sizes="16x16">

  <style>
    body, footer{
      background-color: #343746;
      color: #edf3f7;
    }
    a:hover {
      color: #fff;
    } 
    a, a:link, a:visited {
      color: #ff0d6a;
		}
    code {
			background-color: #282a36;
			white-space: nowrap;
    }
    pre {
			background-color: #282a36;
      border: .0625rem solid #ddd;
    }
  </style>

</head>

<body>
  <div class="container">
    <h2>Welcome to the hosted <a href="https://github.com/aeternity/aesophia">Sophia compiler</a>  for the <a class="special" href="https://aeternity.com">Æternity Blockchain</a></h2>
    <div class="row">
      <div class="col-sm-12 col-md-8 col-lg-5">
        <p>
          The hosted compiler allow you to use different version of the Sophia compiler using the <a
            href="https://github.com/aeternity/aesophia_http">http interface</a>.
        </p>
        <p>
          The available compiler versions are listed below, to understand the difference between versions refer to 
          the official <a href="https://github.com/aeternity/aesophia/blob/master/CHANGELOG.md#changelog">changelog</a>
        </p>

				<ul>
					{{range .Compilers}}
					<li><a href="https://github.com/aeternity/aesophia/releases/tag/{{ .Version }}"><code>{{ .Version }}</code></a>{{if .IsDefault}}(default){{end}}</li>
					{{else}}<p>no compilers avaliable</p>
					{{end}}
        </ul>

        <p>
          To select a specific compiler use the header name <code>{{.Header}}</code> with the value the versions above, 
          if no header is present the default compiler will be used.
        </p>
        <p>
          Here is the curl example to select the compiler version <code>v3.1.0</code>:
        </p>
        <pre>curl -i \
-H "{{.Header}}: v3.1.0" \
-H "Content-Type: application/json" \
https://compiler.aepps.com/api</pre>

      </div>
      
    </div>
    
    
  </div>
  <footer class="sticky">
    <p>v{{.Version}}</p>
  </footer>
</body> 

</html>
`
