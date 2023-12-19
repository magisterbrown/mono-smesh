<script>
import Header from '../Header.svelte';
import Fa from 'svelte-fa/src/fa.svelte'
import { faCheck, faXmark, faBook} from '@fortawesome/free-solid-svg-icons'
import { authenticatedFetch } from '$lib/request.js';
let agents = [
    {id: 0, name: "subm.tar", sucess: true, score: 500},
    {id: 1, name: "agent.tar", sucess: false, score: 200},
    {id: 2, name: "last.tar", sucess: true, score: 600}

];

function uploadAgent(e) {
    authenticatedFetch("/api/leaderboard", "POST", {}, new FormData(e.target));
}
</script>
<Header sel="subm"></Header>
<div class="content">
    <div class="title">Submissions</div>
    <div class="columns">
        <div class="info">
        List of all your active agents. Green ones are active on the leader board.
        Download example agent tar file, unpack it and try to upload your own python agent.
        </div>
        <form class="filed" on:submit|preventDefault={uploadAgent}>
             <input type="file" name="file" id="file" class="inputfile" accept=".tar"/>
             <input type="submit" class="send" value="Upload Agent">
        </form>
    </div>
    <div class="subtitles listed">
        <span class="sucess">Status</span>
        <span class="file">File Name</span>
        <span class="score">Score</span>
        <span class="watch">Games</span>
    </div>
    <div class="list"> 
            {#each agents as agent}
            <div class="agent listed">
                <span class="sucess">
                    {#if agent.sucess}
                    <Fa icon={faCheck} color="#1e8e3e"/>
                    {:else}
                    <Fa icon={faXmark} color="#d93025"/>
                    {/if}
                </span>
                <span class="file">{agent.name}</span>
                <span class="score">{agent.sucess ? agent.score : ""}</span>
                <span class="watch">
                    {#if agent.sucess}
                        <Fa icon={faBook} style="cursor:pointer; font-size: 1.3rem"/>
                    {/if}
                </span>
            </div>
            {/each}
    </div>
</div>

<style>
.title{
    font-weight: 700;
    font-size: 2rem;
}
.subtitles{ 
    color: #333;
    padding: 0.5rem 0 0.5rem 0;
    border-bottom: 1px solid #ccc;
}
.inputfile{
    border: 1px solid #ccc;
    padding: 2rem;
    width: 100%;
    box-sizing: border-box;
    text-align: center;
    border-radius: 1rem;
}
.filed {
    max-width: 20rem;
    margin-left: auto;
}
.filed *{
    display: block;
}

.columns{
    display: flex;
    margin: 1rem 0 1rem  0;
}
.info{
    max-width: 30rem;
    color: #333;
}

.send{
    margin: 0 0 0 auto;
    margin-top: 1rem;
    padding: 0.5rem 1rem 0.5rem 1rem;
    border-radius: 0.7rem;
    border: 1px solid #ccc;
    background: #fff;
    font-weight: 700;
    cursor: pointer;
}
.listed {
    display: flex;
    font-weight: 600;
}

.listed *{
   display: flex;
   align-items: center;
}

.listed .sucess {
    min-width: 5rem;
    text-align: center;
    padding-left: 1.5rem;
}

.agent .sucess {
    font-weight: 800;
    font-size: 3rem;
}

.listed .file {
    flex-grow: 1;
}

.listed .score {
    min-width: 10rem;
    justify-content: center;
}

.listed .watch {
    min-width: 10rem;
    justify-content: center;
}
</style>
