import React from "react";
//@ts-expect-error types not available for react-syntax-highlighter
import SyntaxHighlighter from "react-syntax-highlighter";
import type { MDXComponents } from "mdx/types";

export function useMDXComponents(components: MDXComponents): MDXComponents {
  return {
    p: ({ children }) => <p className="">{children}</p>,
    h2: ({ children }) => <h2 className="font-extrabold my-2 text-lg">{children}</h2>,
    h3: ({ children }) => <h3 className="font-bold my-2">{children}</h3>,
    h4: ({ children }) => <h4 className="font-light my-2">{children}</h4>,
    table: ({ children }) => <table className="my-4 w-full">{children}</table>,
    thead: ({ children }) => <thead className="bg-black text-white border-black border-2">{children}</thead>,
    th: ({ children }) => <th className="text-left p-2">{children}</th>,
    tbody: ({ children }) => <tbody className="border-black border-2">{children}</tbody>,
    tr: ({ children }) => <tr className="text-left">{children}</tr>,
    td: ({ children }) => <td className="text-left p-2">{children}</td>,
    pre: ({ children }) => <pre className="bg-gray-100 p-4 my-4">{children}</pre>,
    code: ({ className, ...properties }) => {
      const match = /language-(\w+)/.exec(className || "");
      return match ? <SyntaxHighlighter language={match[1]} PreTag="div" {...properties} /> : <code className={className} {...properties}></code>;
    },
    BlockQuote: ({ className, type, children }: { className: string; type: string; children: React.ReactNode }) => {
      return (
        <p className={className + " blockquote"}>
          <b>{type}:</b>
          <p>{children}</p>
        </p>
      );
    },
    ...components,
  };
}
