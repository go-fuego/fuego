import Mermaid from "@theme/Mermaid";
import React from "react";

const code = `flowchart TB

subgraph fuego
 direction LR

 subgraph input
	 direction TB
	 req(Request) -- JSON {'a':'value '} --> Deserialization
	 Deserialization -- struct{A:'value '} --> InTransformation
	 InTransformation -- struct{A:'My Value'} --> Validation
 end

 subgraph output
	 direction TB
	 OutTransformation -- struct{B:'My Response!'} --> resp(Response)
 end

 Controller{{Controller}}

 
 
 input -- struct{A:'My Value'} --> Controller -- struct{B:'Response'} --> output
 
 output --> ress(Response)
 
 input -- error --> ErrorHandler
 Controller -- error --> ErrorHandler
 output -- error --> ErrorHandler

 click Controller "/fuego/docs/guides/controllers" "Controllers"
 click Validation "/fuego/docs/guides/validation" "Controllers"
 click InTransformation "/fuego/docs/guides/transformation" "Transformation"
 click OutTransformation "/fuego/docs/guides/transformation" "Transformation"
 click ErrorHandler "/fuego/docs/guides/errors" "Error Handling"
end

Request -- JSON {'a':'value '} --> fuego
fuego -- JSON {'b':'Response!'} --> Response


fuego -- JSON {'error':'Error message', 'code': 4xx} --> Response

`;

export function FlowChart({ selected }) {
  let style = "";
  if (selected && typeof selected === "string") {
    if (selected === "Transformation") {
      style +=
        `style InTransformation stroke:#f33,stroke-width:4px` +
        "\n" +
        `style OutTransformation stroke:#f33,stroke-width:4px`;
    } else {
      style += `style ${selected} stroke:#f33,stroke-width:4px`;
    }
  }

  return <Mermaid value={code + style} />;
}

export default FlowChart;
