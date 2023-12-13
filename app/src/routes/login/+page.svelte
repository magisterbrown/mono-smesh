<script>
    import '../auth.css';
    import { getToken } from '$lib/request.js';
    let username = '';
    let password = '';
    function handleLogin() {
        getToken("/login", username, password);
    }
</script>
<a class="switch" href="/signup">signup</a>
<div class="vert">
    <form class="verification" on:submit|preventDefault={handleLogin}>
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
        <div class="button">
            <input type="submit" class="button" value="Login">
        </div>
    </form>
</div>

