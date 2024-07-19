<svelte:options immutable={false} />

<script lang="ts">
    import { onMount } from "svelte";
    export let data;

    let r = ["loading..."];
    const url: string = '/api/files?';
    async function getReps(): Promise<Array<any>> {
        try {
            const response = await fetch(url  + new URLSearchParams({repo: data.rep}).toString());
            const res = JSON.parse(await response.text()) as string[]
            res.forEach( s => {
                r.push(s)
            })
            r[0] = ""
            console.log(r);
        } catch (error) {
            console.error('Error:', error);
        }
        return []
    }

    onMount(async () => {
        await getReps()
    })

    async function handleSubmit(event: SubmitEvent) {
      const form = event.target as HTMLFormElement;
      const d = new FormData(form);
      window.location.href = '/repos/' + data.rep + '/' + d.get("file") + '/';
    }
</script>

<main>
    <h1>Choose your repository</h1>
    <form on:submit|preventDefault={handleSubmit}>
        <label>
          <span>New file name</span>
          <input name="file" />
        </label>
        <button type="submit">Create new file</button>
    </form>
    <table>
    <tr><p>---------------------</p></tr>
    {#each r as file}
        {#if file !== ''}
            <tr>
                {#if file === "loading..."}
                    <p class="button">{file}</p>
                {:else}
                    <a href="/repos/{data.rep}/{file}/" class="button">{file}</a>
                {/if}
            </tr>
            <tr><p>---------------------</p></tr>
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