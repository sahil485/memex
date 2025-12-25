export namespace main {
	
	export class SearchResult {
	    id: string;
	    path: string;
	    content: string;
	    type: string;
	    title: string;
	    rankingScore: number;
	
	    static createFrom(source: any = {}) {
	        return new SearchResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.path = source["path"];
	        this.content = source["content"];
	        this.type = source["type"];
	        this.title = source["title"];
	        this.rankingScore = source["rankingScore"];
	    }
	}
	export class SearchResponse {
	    hits: SearchResult[];
	    query: string;
	    processingTimeMs: number;
	    estimatedTotalHits: number;
	    error?: string;
	
	    static createFrom(source: any = {}) {
	        return new SearchResponse(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.hits = this.convertValues(source["hits"], SearchResult);
	        this.query = source["query"];
	        this.processingTimeMs = source["processingTimeMs"];
	        this.estimatedTotalHits = source["estimatedTotalHits"];
	        this.error = source["error"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}

}

