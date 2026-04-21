import { defineConfig, defineDocs } from 'fumadocs-mdx/config';
import { metaSchema, pageSchema } from 'fumadocs-core/source/schema';
import lastModified from 'fumadocs-mdx/plugins/last-modified';
import {
  remarkFeedbackBlock,
  type RemarkFeedbackBlockOptions,
} from 'fumadocs-core/mdx-plugins/remark-feedback-block'

const feedbackOptions: RemarkFeedbackBlockOptions = {
  // other options:
};

// You can customise Zod schemas for frontmatter and `meta.json` here
// see https://fumadocs.dev/docs/mdx/collections
export const docs = defineDocs({
  dir: 'content/docs',
  docs: {
    // async:true,
    schema: pageSchema,
    postprocess: {
      includeProcessedMarkdown: true,
    },
  },
  meta: {
    schema: metaSchema,
  },
});

export default defineConfig({
  mdxOptions: {
    // MDX options
    // remarkPlugins: [
    //   [remarkFeedbackBlock, feedbackOptions],
    // ],
  },
  plugins: [
    lastModified()

  ],

});
