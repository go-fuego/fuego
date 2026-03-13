# Release v1.2.3

## What's Changed

* fix: echo adaptor to be compatible with Context params (0ce880c)
* refactor: add dedicated type for engine options (582a572)
* chore(deps): bump the all group across 5 directories with 3 updates (925bdaf)
* Addressing PR comments. (54627fc)
* Added support for openapi3gen schema customizer functions. (df5b061)
* chore(deps): bump the go_modules group across 2 directories with 1 update (08392cc)
* chore(deps): bump gorm.io/gorm (37dfb0e)
* feat: Array of int in golden file example (981d4d4)
* feat: Better error message when passing a pointer to a struct as params (e41678c)
* feat: Fixed error messages for integer overflow in query parameter deserialization (c7d132c)
* feat: Add tests for integer overflow and underflow in query parameter deserialization (85061da)
* feat: Avoid integer overflow when deserializing array of params (ad1a26c)
* refactor: change if/else by switch (9169e82)
* feat: Array of parameters deserialization (b3e3148)
* feat: Add support for more numeral types (4462139)
* Support for float deserialization (42c95bb)
* feat: Support for examples in strongly typed parameters (ee37cde)
* Passes the test for deserializing strongly typed parameters (e1d1e48)
* Fixes types errors (55158ae)
* Support different types of strongly typed query parameters (f9e6d94)
* Used strongly typed parameters in the golden example (df60618)
* Add description to the query parameter (de18d8e)
* Register typed params in OpenAPI spec (0e439e5)
* Strongly typed parameters implementation (7dae127)
* Write and read params from the http response to a struct (428ace2)
* Typed params in context (86ecd39)
* chore(deps): bump the all group across 19 directories with 3 updates (f71b024)
* chore(deps): bump crate-ci/typos from 1.31.2 to 1.32.0 in the all group (5cb4b4e)
* ci: use go-version-file for setting up go in actions (3502f97)
* fix: add struct to log property not found log when parsing structTags (a98f2e4)
* chore(deps): bump the all group across 16 directories with 3 updates (436b2e2)
* chore(deps): bump crate-ci/typos from 1.31.1 to 1.31.2 in the all group (75add4b)
* add more tests & comments (d30a106)
* add tests (bef548c)
* convert base path to fuego format (5860cee)
* check if base path is root (4ce4a6a)
* fix: cast to interface instead of gin.RouterGroup (ba52f00)
* chore: ensure proper context logging (9984311)
* chore: ensure proper context is passed to ErrorHandler (1bf9c18)
* chore: ensure we pass request to Send* functions (91717e0)
* BREAKING: use slog.***Context where applicable (8470afc)
* feat: Parse tags for arrays/slices fields in structs (#500) (af3b109)
* docs: Updated documentation and Youtube link (ae4f30a)
* [Strongly Typed Params] Typed params in context (#406) (d7ffcf7)
* chore: bump golangci-lint version (ef84851)
* Fix: specs endpoint validator and tests (f9f386d)
* Replace TestValidateJsonSpecURL with TestValidateSpecURL (a7814fe)
* feat: Allow non-.json paths for OpenAPI specs (3b23215)
* chore: add tests as well limit internal failure exposure in error responses (29fe7e6)
* fix: add detail/status for queryparam errs and detail for path errs (deac1f0)
* fix: allow setting disable messages from openapi config (b6211e0)
* chore(deps): bump the all group across 2 directories with 2 updates (5f914dc)
* chore(deps): bump golang.org/x/net (667b280)
* chore(deps): bump the go_modules group across 1 directory with 2 updates (dca52a7)
* chore(deps): bump http-proxy-middleware (0b61210)
* chore(deps): bump the go_modules group across 16 directories with 2 updates (467a8ef)
* chore: Updated Go version to 1.24.2 (#484) (212de34)
* chore: Moves sql packages to their own modules (#474) (89ac3b1)
* chore(deps): bump the all group across 14 directories with 1 update (#480) (3c1c9f3)
* feat: Support boolean examples for boolean types in OpenAPI struct tags (#481) (a3806e7)
* Add support for floating point values in example anno (#456) (48365ff)
* fix: default status code for gin/echo (#467) (a95985a)
* feat: Add option.RequestBody (#466) (0111545)
* tools: avoid redundant fmt from golangci-lint + remove exclusions (#473) (217ca9a)
* Golangci lint v2 (#464) (c30da83)
* chore(deps): bump the all group across 17 directories with 9 updates (#472) (06686eb)
* chore(deps): bump the all group across 2 directories with 5 updates (#471) (8606f36)
* chore(deps): bump the all group with 2 updates (#470) (2efedf7)
* chore: Dependabot configuration (#469) (47a33e0)
* chore: Codeowners file (#462) (8b05a6e)
* chore(deps): bump golang.org/x/net (#461) (a3006ea)
* chore(deps): bump the npm_and_yarn group across 1 directory with 3 updates (#460) (0105b00)
* chore(deps): bump the go_modules group across 19 directories with 1 update (#458) (3859fd6)
* chore(deps): bump the npm_and_yarn group across 1 directory with 3 updates (#452) (2cdb235)
* chore: Bump deps jwt (#457) (7c0db81)
* feat: Add missing fuegogin methods (#447) (3424a3d)
* feat: Allow use of gin path format (#449) (545b0c0)
* FIX: duplicate base path when using gin groups (#445) (0337182)
* fix: skip calling Flush on responseWriter when it is not implemented (#443) (59b4ec2)
* chore: Run modernize (3ba7dd2)
* feat: add Server option of WithStripTrailingSlash (#432) (4db8a77)
* Accepts list as input in the body of the request (#440) (2e22ebf)
* chore: update documentation to provide working example (#438) (1a2e09b)
* feat: Better query params logging when unexpected in OpenAPI spec (f331854)
* feat: Unify HTTP errors and improve error messages (#431) (a296bd1)
* style: Rewords the panic error when a path parameter is declared in OpenAPI but not on the route (151afe9)
* feat: Support http.ResponseWriter options (Flusher, Pusher, Hijacker) (#430) (868535f)
* feat: Added options methods (#429) (9a7a85c)
* extra: map SQL and SQLite errors to Fuego HTTP errors (#427) (1932a81)
* docs: fix OpenApi Specification example code (#426) (eaf297d)
* feat: option.DefaultResponse (c8a2d52)
* feat: add option.DefaultResponse (#425) (58fcd26)
* chore: improve test readability & add testifylint (#422) (93f50ab)
* docs: Middlewares guide update (c5cdb3c)
* docs: Updated the error handling guide (59d6d56)
* docs: Updated documentation website URL (559f36c)
* docs: Updated Docusaurus baseUrl (dd5374b)
* refactor: Cleanup redundant path/query param funcs (#421) (26a2f85)
* Adds unconvert linter (#423) (6585949)
* Updated documentation (61ab99e)
* example: Generate openapi spec without starting server using tools pattern (#409) (276e81b)
* Update engine.go (404174b)
* option to disable default open api server (708bcd2)
* feat: Log error even if not fuego.HTTPError (#403) (f73f3f2)
* feat: expose HandleHTTPError for use WithErrorHandler option (c21f647)
* feat: expose HandleHTTPError for use WithErrorHandler option (9a4a4ea)

## Contributors

* @Benjamin Meyer
* @charlie
* @Conner Murphy
* @Cristian
* @dependabot[bot]
* @Dylan Hitt
* @dylanhitt
* @Evilenzo
* @Ewen Quimerc'h
* @EwenQuim
* @Hendrik Ohrdes
* @Jonathan Witchard
* @Josh
* @Juan Miguel  Arias Mejias
* @Majo Richter
* @Malcolm Rebughini
* @Oleg Kunitsyn
* @Philip Jöbstl
* @root
* @serchemach
* @sirmackan
* @Stone Olguin
