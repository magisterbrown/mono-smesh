<script>
    import '../auth.css';
    import { getToken } from '$lib/request.js';
    let username = '';
    let password = '';
    let repeat = '';
    function handleSignup() {
        let error = false;
        error |= (username.length<2);
        error |= (password.length<2);

        if(error){
            alert("badfomr")
        }
        getToken("/signup", username, password);
    }
</script>
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
        <input type="username" class="textf" id="username" bind:value={username}>
        <label for="password">Password:</label>
        <input type="password" class="textf" id="password" bind:value={password}>
        <label for="confirm">Confirm password:</label>
        <input type="password" class="textf" id="confirm" bind:value={repeat}>
        <div class="button">
            <input type="submit" class="button" value="Sign Up">
        </div>
    </form>
</div>

