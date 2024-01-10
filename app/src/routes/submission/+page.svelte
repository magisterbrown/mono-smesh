<script>
import Header from '../Header.svelte';
import Matches from '../Matches.svelte';
import Response from './Response.svelte';
import Fa from 'svelte-fa/src/fa.svelte'
import { faCheck, faXmark, faBook} from '@fortawesome/free-solid-svg-icons'
import { authenticatedFetch } from '$lib/request.js';
import { onMount } from "svelte";
let agents = [];
let showAgent;
let respUpload;
let owner;
let waitingUpload=false;

onMount(() => {
    authenticatedFetch("/api/whoami", {method: "GET"}).then(resp => {
        resp.text().then(name => {
            owner=name
            authenticatedFetch("/api/submissions?user_name="+name, {method: "GET"}).then(resp => {
                     resp.json().then(body => {agents = body
                     });
            });
        });
    });
});

function uploadAgent(e) {
    waitingUpload=true;
    authenticatedFetch("/api/leaderboard", {method: "POST", body: new FormData(e.target)}).then(resp => {
       if(resp.ok) {
            resp.json().then(agent => {
                agents.unshift(agent) 
                agents = agents
            })
       } else {
            resp.text().then(errmsg => {respUpload=errmsg})
       }
       waitingUpload=false;
    });
}

const intervals = [
  { label: 'year', seconds: 31536000 },
  { label: 'month', seconds: 2592000 },
  { label: 'day', seconds: 86400 },
  { label: 'hour', seconds: 3600 },
  { label: 'minute', seconds: 60 },
  { label: 'second', seconds: 1 }
];

function timeSince(date) {
  const seconds = Math.floor((Date.now() - Date.parse(date)) / 1000)+2;
  const interval = intervals.find(i => i.seconds < seconds);
  const count = Math.floor(seconds / interval.seconds);
  return `${count} ${interval.label}${count !== 1 ? 's' : ''} ago`;
}

$: console.log(agents);
</script>

<Response bind:respUpload></Response>
<Matches bind:showAgent bind:owner> </Matches>

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
             <div class="sendwrap">
                <input type="submit" class="send" value="Upload Agent" style="display: {waitingUpload ? 'none' : 'block'};">
             </div>
        </form>
    </div>
    <div class="subtitles listed">
        <span class="sucess">Status</span>
        <span class="file">File Name</span>
        <span class="score">Score</span>
        <span class="watch">Games</span>
    </div>
    <div class="list"> 
            {#if agents.length>0}
            {#each agents as agent}
            <div class="agent listed">
                <span class="sucess">
                    {#if !agent.Broken}
                    <Fa icon={faCheck} color="#1e8e3e"/>
                    {:else}
                    <Fa icon={faXmark} color="#d93025"/>
                    {/if}
                </span>
                <span class="file">{agent.FileName}<span class="timeago"> ({timeSince(agent.CreatedAt)})</span></span>
                <span class="score">{!agent.Broken ? agent.Raiting.toFixed(0) : ""}</span>
                <a class="watch" on:click|stopPropagation={()=>{showAgent=agent.Id}}>
                    {#if !agent.Broken}
                        <Fa icon={faBook} style="cursor:pointer; font-size: 1.3rem"/>
                    {/if}
                </a>
            </div>
            {/each}
            {/if}
    </div>
</div>

<style>
.timeago{
    color: #333;
    padding-left: 0.5rem;
}
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
.sendwrap{
    min-height: 2.5rem;
    margin-top: 0.8rem;
}
.send{
    margin-top: 1rem;
    border-radius: 0.7rem;
    border: 1px solid #ccc;
    background: #fff;
    font-weight: 700;
    cursor: pointer;
    padding: 0.5rem 1rem 0.5rem 1rem;
    margin: 0 0 0 auto;
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
