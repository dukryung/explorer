<script>
    export let page = 0;
    export let pageLimit = 0;
    export let maxItem = 0;

    let firstPage = false;
    let previousPage = false;
    let nextPage = false;
    let lastPage = false;

    function onFirst() {
      page = 0;
    }

    function onLast() {
      page = maxPage - 1;
    }

    function onPrevious() {
      page = page - 1;
    }

    function onNext() {
      page = page + 1;
    }

    function dim(state) {
      return state === false ? 'PaginatorDim' : '';
    }

    $: maxPage = Math.ceil(maxItem / pageLimit);
    $: if (page || maxPage) {
      if (page === 0) {
        firstPage = false;
        previousPage = false;
      } else {
        firstPage = true;
        previousPage = true;
      }

      if (maxPage > 0 && page !== maxPage - 1) {
        nextPage = true;
        lastPage = true;
      } else {
        nextPage = false;
        lastPage = false;
      }
    }
</script>

<div class="Paginator BlockListpaginator">
    <ul>
        <li class="PaginatorBtn PaginatorFirstPage {dim(firstPage)}" on:click={onFirst}>
            <a><img src="/images/firstPage.svg" class=""></a>
        </li>
        <li class="PaginatorBtn PaginatorPrevBtn {dim(previousPage)}" on:click={onPrevious}>
            <a><img src="/images/PrevPage.svg" class=""></a>
        </li>
        <li class="PaginatorpageNumber">
            <p class="PaginatorCurrentNum">{page + 1} / {maxPage === 0 ? 1 : maxPage}</p>
        </li>
        <li class="PaginatorBtn PaginatorNextBtn {dim(nextPage)}" on:click={onNext}>
            <a><img src="/images/NextPage.svg" class=""></a>
        </li>
        <li class="PaginatorBtn PaginatorLastPageBtn {dim(lastPage)}" on:click={onLast}>
            <a><img src="/images/lastPage.svg" class=""></a>
        </li>
    </ul>
</div>