<script>
import Header from '../Header.svelte';
let agents = [
    {id: 0, sucess: true, score: 500},
    {id: 1, sucess: false, score: 200},
    {id: 2, sucess: true, score: 600}

];
</script>
<Header sel="subm"></Header>
<div class="content">
    <div class="title">Submissions</div>
    <div class="subtitles">
        <span class="sucess">Status</span>
        <span class="score">Score</span>
        <span class="watch">Games</span>
    </div>
    <div class="list"> 
            {#each agents as agent}
            <div class="agent">
                <span class="sucess"></span>
                <span class="score">{agent.score}</span>
                <span class="watch"></span>
            </div>
            {/each}
    </div>
</div>

<style>

</style>
