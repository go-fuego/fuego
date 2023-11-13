package op

// HTML is a marker type used to differentiate between a string response and an HTML response.
// To use templating, use [Ctx.Render].
type HTML string
