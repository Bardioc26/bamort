## TODO (Remaining)
* 1. in RenderPageWithContinuations use a switch to select a templatename and the ncall a paginator function that is designed explicitly for this template page with the current content. This will not work when we add more and different templates. We nned a solution where the code can find out which data is present on the template page and if these datastructures are to be continued on a continuation page or on the next regular page in row. it needs to find out how many elements can be rendered into the page or it's continuation page. It needs to handle if there are multiple datastructures present on one page.
for this to be accomplished we need some kind of template metadata or inline datastructure(Block) configuration.
I want to change the templateName from page1_stats.html to page_1.html and accordingly the continuation page from page1.2_stats.html to page_1.2.html this makes clear that a page is not tight to one or another datastructure.
Page 4 has a complex container-based layout we leave this out for the moment.
Plan a refacturation of PaginateSkills, PaginatePage2PlayLists and PaginateSpells into one unified function. 



### Later
* continuation of lists does not work as expected but good enough for a first shot
  * generalize handling so that only on set of functions can handle ALL kinds of templates. Needs massive refactoring

* currently the template fetched for rendering is set to Default_A4_Quer
* remove inline css as far as possible
* make pdf download popup an own view
* func CleanupExportTemp move maxAge := 7 * 24 * time.Hour definition to Config struct n config.go