import {Typography} from '@mui/material';
import {CommonProps} from '@mui/material/OverridableComponent';

export function Explain(props: CommonProps) {
  return (
    <>
      <Typography variant="h4" sx={{mt: 2}}>Digital Dexterity</Typography>
      <Typography variant="h6">Willingness To Learn</Typography>
      <Typography variant="subtitle2">How willing are they to learn about a new technology/technique that may or may not effect their life</Typography>
      <Typography variant="h6">Self Sufficient Learning</Typography>
      <Typography variant="subtitle2">
        Self-sufficient people have the ability and the desire to live life on their own terms, to determine their own course, to make their own decisions, and not have their life choices made by others.
        <br />Are they able to choose what to learn or do they need someone to tell them
      </Typography>
      <Typography variant="h6">Improving Capability</Typography>
      <Typography variant="subtitle2">
        Do they actively seek out improving existing capabilities at their place of work/home and/or create new capabilities.
      </Typography>
      <Typography variant="h6">Innovative Thinking</Typography>
      <Typography variant="subtitle2">
        Do you think of them as someone who comes up with new problems or solutions that no one else may have thought of yet?
      </Typography>
      <Typography variant="h6">Growth Mindset</Typography>
      <Typography variant="subtitle2">
        Are they able to identify areas of growth needed in yourself/others?
      </Typography>
      <Typography variant="h6">Awareness Of Self Efficacy</Typography>
      <Typography variant="subtitle2">
        Self-efficacy refers to an individual&apos;s belief in his or her capacity to execute behaviors necessary to produce specific performance attainments (Bandura, 1977, 1986, 1997). Self-efficacy reflects confidence in the ability to exert control over one&apos;s own motivation, behavior, and social environment.
      </Typography>
      <Typography variant="h6">Applying What They Learn</Typography>
      <Typography variant="subtitle2">
        Do they apply in practicality something that they have just been taught, e.g. if they learned cloud concepts recently can they actually use cloud concepts or do they still reluct to.
      </Typography>
      <Typography variant="h6">Adaptability</Typography>
      <Typography variant="subtitle2">
        Can they change how they act/behave and/or can they improve based on circumstances that are thrown at them.
      </Typography>
    </>
  );
}
