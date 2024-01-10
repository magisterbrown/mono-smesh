<script>
    export let respUpload=undefined;
</script>
{#if respUpload !== undefined}
<div class="popup error" >
    <div class="header">
        Uploading Error: 
        <button class="exit" on:click|preventDefault={() => {respUpload=undefined}}>
            <div class="bar" style="transform: rotate(45deg)"></div>
            <div class="bar"style="transform: rotate(-45deg)"></div>
        </button>
        <div>{respUpload}</div>
    </div>
</div>
{/if}
<style>
.error{
    background-color: #dc3545;
}
.correct{
    background-color: #198754;
}
.popup{
    font-weight: 600;
    color: #fff;
    position: fixed;
    width: 35rem;
    max-width: 100%;
    left: 50%;
    top: 7rem;
    transform: translate(-50%, 0);
    padding: 0.7rem;
}
.exit{
    float: right;
    position: relative;
    cursor: pointer;
    width: 2.5rem;
    height: 2.5rem;
    background: none;
	color: inherit;
	border: none;
	padding: 0;
	font: inherit;
	cursor: pointer;
	outline: inherit;
}

.bar {
    top:0%; 
    right:50%;
    transform: translate(-50%, -50%);
    position: absolute;
    width: 0.3rem;
    height: 2.2rem;
    background-color: #fff;
}
</style>
