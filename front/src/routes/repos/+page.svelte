<svelte:options immutable={false} />

<script lang="ts">
    import { onMount } from "svelte";

    let r = ["loading or bad github token, check logs"];
    const url: string = '/api/rep';
    async function getReps(): Promise<Array<any>> {
        try {
            const response = await fetch(url);
            const res = JSON.parse(await response.text()) as string[]
            res.forEach( s => {
                r.push(s)
            })
            r[0] = ""    // wtf why does the array update don't work, it's fucking insane, I've spent on this thing few hours
            console.log(r);
        } catch (error) {
            console.error('Error:', error);
        }
        return []
    }

    onMount(async () => {
        await getReps()
    })
</script>

<main>
    <h1>Choose your repository</h1>
    <table>
    {#each r as rep}
        {#if rep !== ''}
            <tr>
                <a href="/{rep}/new" class="button">{rep}</a>
            </tr>
        {/if}
    {/each}
    </table>
</main>

<style>
    main {
        text-align: center;
        padding: 1em;
        max-width: 240px;
        margin: 0 auto;
    }

    h1 {
        color: #ff3e00;
        text-transform: uppercase;
        font-size: 4em;
        font-weight: 100;
    }

    table {
        color: #ff3e00;
        margin-left: auto;
        margin-right: auto;
        font-size: 1.5em;
        font-weight: 100;
    }

    @media (min-width: 640px) {
        main {
            max-width: none;
        }
    }
</style>