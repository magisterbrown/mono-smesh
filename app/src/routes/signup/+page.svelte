<script>

    import '../auth.css';
    import { auth_url} from '../config.js';
    let accInfo = {
        username:'',
        password:'',
    }
    let repeat = '';
    function handleSignup() {
        let error = false;
        error |= (accInfo.username.length<2);
        error |= (accInfo.password.length<2);

        if(error){
            alert("badfomr")
        }
        fetch(auth_url+"/signup", { method: "POST", mode: "cors", headers: {"Content-Type": "application/json"},body: JSON.stringify(accInfo)}).then(data => {
           console.log(data); 
        });
        console.log(auth_url);

    }
</script>
<body>
    <a class="switch" href="/login">login</a>
    <div class="vert">
        <form class="verification" on:submit|preventDefault={handleSignup}>
            <div class="logo">
                <div class="epoch">
                    <div id="cir1"></div>
                    <div id="cir2"></div>
                </div>
            </div>
            <label for="username">User Name:</label>
            <input type="username" class="textf" id="username" bind:value={accInfo.username}>
            <label for="password">Password:</label>
            <input type="password" class="textf" id="password" bind:value={accInfo.password}>
            <label for="confirm">Confirm password:</label>
            <input type="password" class="textf" id="confirm" bind:value={repeat}>
            <div class="button">
                <input type="submit" class="button" value="Sign Up">
            </div>
        </form>
    </div>

</body>
