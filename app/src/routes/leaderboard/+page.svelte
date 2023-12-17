<script>
import Header from '../Header.svelte';
import Fa from 'svelte-fa/src/fa.svelte'
import { faBook} from '@fortawesome/free-solid-svg-icons'
import { authenticatedFetch } from '$lib/request.js'
import { onMount } from "svelte";

onMount( () => {
    authenticatedFetch("/api/leaderboard");
});
let teams = [
    { id: 1, name: "BigTeam", score: 1200.3, agents: 5},
    { id: 2, name: "BigTeamU", score: 1101.1, agents: 15},
    { id: 3, name: "LuTeam", score: 45.3, agents: 3695}
]
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
            {#each teams as team, i}
                <div class="listed leader">
                    <span class="rank">{i+1}</span>
                    <span class="team">{team.name}</span>
                    <span class="score">{team.score}</span>
                    <span class="agents">{team.agents}  <Fa icon={faBook} style="cursor:pointer; font-size: 1.0rem"/>
</span>
                </div>
            {/each}
            
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
