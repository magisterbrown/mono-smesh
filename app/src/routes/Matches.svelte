<script>
    import {onMount} from 'svelte';
    import { authenticatedFetch } from '$lib/request.js'
    import { backend_url } from '$lib/config.js'
    import Fa from 'svelte-fa/src/fa.svelte'
    import { faPlay } from '@fortawesome/free-solid-svg-icons'
    let matches = [];
    export let showAgent = undefined;;
    export let owner;
    let gameId = undefined;
    $: if( showAgent !== undefined) {
        authenticatedFetch("/api/matches?id="+showAgent, {method: "GET"}).then(resp => {
            resp.json().then(body => {body===null ? matches=[] : matches=body})
        });
    }
    function hide() {
        showAgent=undefined
        gameId=undefined
        matches=[]
        clearInterval(intervalId);
    }
    let intervalId; 
   $: if(gameId !== undefined){
        authenticatedFetch("/api/matches/recording?id="+gameId).then(resp => {
            resp.json().then(body => {
                let cst = 0;
                intervalId = setInterval(function() {
                    let state = body[cst%body.length];
                    cst++;
                    render(state)
                    //Module.ccall("c_add", "number", ["string"],[state])
                }, 1000);

            })
        });
    }
    function line(ctx, x0, y0, x1, y1){
        ctx.beginPath();
        ctx.moveTo(x0, y0);
        ctx.lineTo(x1, y1);
        ctx.stroke(); 
     
    }
    function render(board_state) {
        const canvas = document.getElementsByClassName("video")[0];
        const ctx = canvas.getContext("2d");
        ctx.clearRect(0, 0, canvas.width, canvas.height);
        const offset = 40;
        const spacing = 80;
        var posx = offset;
        var posy = offset;
        var fslash = false;
        var toprow = true;
        var leftcol = true;
        var figures = [];
        const board = board_state.state.split(' ')[0];
        for (const spot of board) {
            if(spot==='/'){
                posx=offset;
                posy+=spacing;
                toprow=false;
                leftcol=true;
            }
            else{
                var cells = parseInt(spot);
                if(isNaN(cells)){
                    cells=1;
                    figures.push({posx: posx, posy: posy, color: spot})
                }
                for(var i=0; i<cells; i++){
                    if(!leftcol)
                        line(ctx, posx, posy, posx-spacing, posy);
                    if(!toprow)
                        line(ctx, posx, posy, posx, posy-spacing);
                    if(!leftcol && !toprow) {
                        if(fslash)
                            line(ctx, posx-spacing, posy, posx, posy-spacing);
                        else
                            line(ctx, posx, posy, posx-spacing, posy-spacing);
                        }
                    fslash = !fslash;
                    posx+=spacing;
                    leftcol=false;
                }
            }

        }
        for(const fig of figures) {
            ctx.beginPath();
            ctx.arc(fig.posx, fig.posy, spacing/4, 0, 2 * Math.PI);
            ctx.fillStyle = fig.color === "W" ? "#c8c8c8" : "#000";
            ctx.fill();
        }
    }
</script>
{#if showAgent !== undefined}
<div class="background" on:click|stopPropagation={hide}>
    <div class="content" on:click|stopPropagation={()=>{}}>
        {#if gameId !== undefined}
            <canvas class="video" height="400" width="800">
            </canvas>
        {:else}
            <div class="tops">
                <div class="header">
                    Submission {showAgent} 
                </div>
                <div class="list">
                {#each matches as match, i}
                    <div class="match" on:click|stopPropagation={() => {gameId=match.Id }}>
                        <div class="players">
                            {#each match.Seating as player, ii}
                                {ii===0 ? "" : " vs"}
                                <span style:font-weight={player.UserName===owner ? 600 : 500}>
                                [{player.Status}]
                                {player.UserName} 
                                ({player.Change>0 ? "+" : ""}{player.Change.toFixed(0)})   
                                </span>
                            {/each}
                        </div>
                        <div class="play"><Fa icon={faPlay}/></div>
                    </div>
                    {/each}
                </div>
            </div>
        {/if}
        <div class="footer">
             <button type="button" class="close" on:click|stopPropagation={hide}>Close</button> 
        </div>
    </div>
</div>
{/if}


<style>
.list{
    max-height: 25rem;
    overflow-y: auto;
}
.video{
    height: 23rem;
    margin: 1.5rem;
}
.background{
    display: grid;
    place-items: center;
    position: fixed;
    width: 100%;
    height: 100%;
    background-color: #00000068;
    z-index: 1;
}
.content{
    width: 45rem;
    max-width: 100%;
    position: relative;
    background-color: #fff;
    border-radius: 2rem;
    padding: 0;
}
.tops{
    margin: 2.5rem 2.5rem 0.5rem 2.5rem;
}
.header{
    font-weight: 600;
    font-size: 2rem;
    margin-bottom: 1.5rem;
}
.match{
    height: 3.5rem;
    display: flex;
    align-items: center;
    padding-left: 0.5rem;
    cursor: pointer;
}
.match:hover{
    background-color: #eee;
}

.players{
    width: 100%;
}

.play{
    margin: 0 1rem 0 1rem;
}

.footer{
    height: 3.5rem;
    border-top: 1px solid #ccc;
    position: relative;
}

.close{
    position: absolute;
    top: 50%;
    right: 2rem;
    transform: translate(0, -50%);
    border-radius: 0.7rem;
    border: 1px solid #ccc;
    background: #fff;
    font-weight: 700;
    cursor: pointer;
    padding: 0.5rem 1rem 0.5rem 1rem;

}

</style>

