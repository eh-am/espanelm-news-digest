import scrapePage from './scrapePage';
import scrapeRSS from './scrapeRSS';
import { map, mergeMap, reduce } from 'rxjs/operators';
import { Observable, of, EMPTY } from 'rxjs';

// we need to infer the type later
type Unpack<T> = T extends Observable<infer U> ? U : never;
type Article = Unpack<ReturnType<typeof scrapePage>>;

export function run(url: string): Observable<any> {
  return scrapeRSS(url).pipe(
    mergeMap(a => scrapePage(a.link)),
    mergeMap(a => (a.links.es ? of(a) : EMPTY)),
    mergeMap(a => {
      return scrapePage(a.links.es).pipe(map(r => mountObject(a, r)));
    }),
    // we are only interested in the last values
    reduce((prev, curr) => [...prev, curr], [] as ReturnType<
      typeof mountObject
    >[])
  );
}

function mountObject(
  ptbr: Article,
  es: Article
): {
  ptbr: Article;
  es: Article;
} {
  return {
    ptbr,
    es
  };
}
