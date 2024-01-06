<script>
    import {onMount} from 'svelte';
    import { authenticatedFetch } from '$lib/request.js'
    export let submissionId=2;
    let matches = [];
    onMount(() => {
        authenticatedFetch("/api/matches?id="+submissionId, {method: "GET"}).then(resp => {
            resp.json().then(body => {matches=body})
        });
    });
</script>

{#if submissionId}
<script>
 console.log('sadas');
</script>
<div class="background" on:click|stopPropagation={()=>{show=false}}>
    <div class="content"  on:click|stopPropagation={()=>{}}>
        <div class="header">
            Submission {submissionId} 
        </div>
        {#each matches as match, i}
        <div>{match.Id}</div>
        {/each}
    </div>
</div>
{/if}


<style>
.background{
    display: grid;
    place-items: center;
    position: absolute;
    width: 100%;
    height: 100%;
    background-color: #00000068;
    z-index: 1;
}
.content{
    width: 45rem;
    height: 30rem;
    max-width: 100%;
    position: relative;
    background-color: #fff;
    border-radius: 2rem;
    padding: 2.5rem;
}
.header{
    font-weight: 600;
    font-size: 2rem;
    margin-bottom: 1rem;
}
</style>

