<script>
import Header from '../Header.svelte';
import Fa from 'svelte-fa/src/fa.svelte'
import { faBook} from '@fortawesome/free-solid-svg-icons'
import { authenticatedFetch } from '$lib/request.js'
import { onMount } from "svelte";

let teams = [];
let my_name = ""

onMount( () => {
    authenticatedFetch("/api/whoami", {method: "GET"}).then(resp => {
        resp.text().then(name => {my_name = name})
    });
    authenticatedFetch("/api/leaderboard", {method: "GET"}).then(resp => {
        resp.json().then(body => {teams = body})
    });
});

</script>
<Header sel="leader"></Header>
<div class="content">
    <div class="table">
        <div class="titles listed">
            <span class="rank">#</span> 
            <span class="team">Team</span>
            <span class="score">Score</span>
            <span class="agents">Agents</span>
        </div>
        <div >
            {#if teams.length}
                {#each teams as team, i}
                    <div class="listed leader" style:font-weight={team.Name==my_name ? 'bold' : 'normal'}>
                        <span class="rank">{i+1}</span>
                        <span class="team">{team.Name}</span>
                        <span class="score">{team.Raiting.toFixed(0)}</span>
                        <span class="agents">{team.Agents}  <Fa icon={faBook} style="cursor:pointer; font-size: 1.0rem"/></span>
                    </div>
                {/each}
            {/if}
            
        </div>
    </div>
</div>

<style>
.table{
   width: 100%;
   max-width: 65rem;
   margin: 0 auto; 
   display: block;
   border: 1px solid #ccc;
   border-radius: 1rem;
}
.titles{
    font-weight: 700;
}

.listed{
    padding: 0.4rem;
    display: flex;
}

.listed * {
    margin-right: 1%;
}

.leader{
    font-weight: 400;
    border-top: 1px solid #ccc;
}

.rank{
    width: 15%;
    text-align: center;
}
.team{
    width: 55%;
}
.score{
    text-align: right;
    width: 15%;
}
.agents{
    text-align: right;
    width: 15%;
}
</style>
