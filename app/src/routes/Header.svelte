<script>
    import { onMount } from 'svelte';
    export let sel;
    let leader = () => {window.location.href = "/leaderboard"}
    let subm = () => {window.location.href = "/submission"}
</script>
<div class="header">
    <div class="lists">
        <div class="buttons">
            <span class:selected={sel=="leader"} on:click={leader}>Leaderboard</span>
            <span class:selected={sel=="subm"} on:click={subm}>Submission</span>
        </div>
    </div>
</div>

<style>
.header{
    width: 100%;
    border-bottom: 1px solid #ccc;
    margin-bottom: 3rem;
    color: #333;
}
.lists{
    max-width: 65rem;
    margin: 0 auto;
    display: flex;
}
.buttons{
    display: flex;
    height: 4rem;
}

.buttons * {
  line-height: 4rem;
  padding: 0 0.8rem 0 0.8rem;
  cursor: pointer;
}

.selected {
    color: #000;
    position: relative;
}
.selected::after {
    content: "";
    position: absolute;
    bottom: 0;
    left: 0;
    background: #000;
    height: 0.20rem;
    width: 100%;
    border-radius: 1rem 1rem 0 0;
}
</style>
