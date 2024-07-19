<script lang="ts">
    let err = "";

    async function handleSubmit(event: SubmitEvent) {
      const form = event.target as HTMLFormElement;
      const data = new FormData(form);

      console.log(data);
      console.log(data.get("tkn"));
  
      let url = '/api/set/?';
      let tknVal = data.get("tkn")?.toString();
      if (tknVal === undefined) {
        tknVal = ""
      }
      const response = await fetch(url  + new URLSearchParams({tkn: tknVal}).toString(), {
        method: 'Get',
      });
      console.log(response.status);
      if (response.status !== 200) {
        err = await response.text();
      } else {
        window.location.href = '/repos';
      }
    }
</script>
  
<main>
<h1>Hey, set up your gihub token!</h1>
<form on:submit|preventDefault={handleSubmit}>
  <label>
    <span>Github token</span>
    <input type="password" name="tkn" />
  </label>
  <button type="submit">Set token</button>
  {#if err !== ''}
    <span style="color: #ff3e00;">{err}</span>
  {/if}
</form>

</main>
  

<style>
    main {
        text-align: center;
        padding: 1em;
        max-width: 240px;
        margin: 0 auto;
    }

    table {
        color: #ff3e00;
        margin-left: auto;
        margin-right: auto;
        font-size: 2em;
        font-weight: 100;
    }

    h1 {
        color: #ff3e00;
        text-transform: uppercase;
        font-size: 4em;
        font-weight: 100;
    }

    @media (min-width: 640px) {
        main {
            max-width: none;
        }
    }
</style>