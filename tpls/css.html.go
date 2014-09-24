package tpls

func init() {
	registerTemplate("css.html", `
<style type="text/css">
        *{border:0; padding:0; margin:0;}
        li {
            list-style: none;
        }

        body {
            background: #e2e2e2
        }
        a {
            text-decoration: none;
            color: #000;
        }

        a:hover {
            background: #ddd;
        }
        .active {
            background: #ddd;
        }
        #box {
            width: 600px;
            height: 500px;
            border-bottom: 1px solid #e2e2e2;
            margin:20px  auto;
            padding: 4px 10px;
            box-shadow: 0px 2px 3px;
            border-radius:3px;
            background: #fff;
        }

        h1 {
            font-size: 18px;
            line-height: 30px;
            color:#999;
            text-indent: 10px;
        }

        #box_heard {
            padding: 1px;
            font-size: 14px;
            line-height: 30px;
            text-align: left;
            border-bottom: 1px solid #e2e2e2;
        }

        #box_body {
            padding: 10px;
            font-size: 12px;
            line-height: 30px;

        }

        .btn {
            width: 40px;
            height: 50px;
            text-align: center;
            display: inline-block;
            padding: 0 15px;
            border-radius: 3px;
            line-height: 50px;
            border:1px solid #D9D9D9;
            margin: 0 30px 0 0;
        }

        #nav {
            width: auto;
            margin: 20px auto;
            text-align: center;
        }

        #text {
            padding: 10px 0 0 5px;
            color: #999;
            line-height: 24px;
        }

        #text p {
            padding-left: 8px;
        }

        .sl {
            border-radius: 3px;
            padding: 5px;
            border: 1px solid #ccc;
            width: 300px;
        }

        #form {
            margin:20px 0 10px 0;
        }

        #sub {
            width:50px;
            background-color: #f9f9f9;
            border: 1px solid rgba(60,60,70,0.3);
            color: #333;
            text-shadow:0px 1px 0px #fff;
            font-weight: #fffbold;
            padding: 3px 4px;
            border-radius: 3px;
            margin-top: 10px;
        }
    </style>
`)
}
