import {Box, Button, Dialog, DialogContent, DialogTitle, Divider, Alert, Slider, TextField, Typography} from '@mui/material';
import {HTTPError} from 'ky';
import {useState} from 'react';
import {useNavigate} from 'react-router-dom';
import {useNewAssessment} from '../../api/assessments';
import {useUsername} from '../../api/profile';

export default function AssessmentForm(props: { self?: boolean }) {
  const [who, setWho] = useState('');
  const [willingnessToLearn, setWillingnessToLearn] = useState(50);
  const [selfSufficientLearning, setSelfSufficientLearning] = useState(50);
  const [improvingCapability, setImprovingCapability] = useState(50);
  const [innovativeThinking, setInnovativeThinking] = useState(50);
  const [growthMindset, setGrowthMindset] = useState(50);
  const [awarenessOfSelfEfficacy, setAwarenessOfSelfEfficacy] = useState(50);
  const [applyingWhatTheyLearn, setApplyingWhatTheyLearn] = useState(50);
  const [adaptability, setAdaptability] = useState(50);
  const [error, setError] = useState('');
  const username = useUsername();
  const {selfAssessment, externalAssessment} = useNewAssessment(username);
  const navigate = useNavigate();
  async function doAssessmentCreation() {
    if (!props.self && who === '') {
      setError('You need to select a user');
    } else {
      const payload = {
        WillingnessToLearn: willingnessToLearn,
        SelfSufficientLearning: selfSufficientLearning,
        ImprovingCapability: improvingCapability,
        InnovativeThinking: innovativeThinking,
        GrowthMindset: growthMindset,
        AwarenessOfSelfEfficacy: awarenessOfSelfEfficacy,
        ApplyingWhatTheyLearn: applyingWhatTheyLearn,
        Adaptability: adaptability,
      };
      if (props.self) {
        selfAssessment.mutate(payload, {onError: async (e) => {
          const reason = (await (e as HTTPError).response.json()).error;
          setError(reason);
        }});
      } else {
        (payload as any).assessing = who;
        externalAssessment.mutate(payload as any, {onError: async (e) => {
          const reason = (await (e as HTTPError).response.json()).error;
          setError(reason);
        }});
      }
    }
  }
  return (
    <>
      <Dialog open={Boolean(selfAssessment.isSuccess || externalAssessment.isSuccess)} onClose={() => navigate('/')}>
        <DialogTitle>Successfully Created Assessment</DialogTitle>
        <DialogContent>
          <Alert severity="success">You have successfully completed assessment</Alert>
          <Button onClick={() => navigate('/')}>OK</Button>
        </DialogContent>
      </Dialog>
      <Dialog open={Boolean(error)} onClose={() => setError('')}>
        <DialogTitle>Failed To Create Assessment</DialogTitle>
        <DialogContent>{error}</DialogContent>
      </Dialog>
      <Box component="form" maxWidth={600}>
        <Typography variant="h5">Digital Dexterity Assessment</Typography>
        <Typography variant="subtitle1">Assessing - {props.self ? 'Self' : 'Someone Else'}</Typography>
        {!props.self && <>
          <Divider sx={{my: 2}} />
          <TextField label="Who Are You Assessing?" value={who} onChange={(e) => setWho(e.target.value)} fullWidth />
        </>}
        <Divider sx={{my: 2}} />
        <Typography variant="h6">Willingness To Learn</Typography>
        <Typography variant="subtitle2">How willing are they to learn about a new technology/technique that may or may not effect their life</Typography>
        <Slider value={willingnessToLearn} onChange={(_, val) => setWillingnessToLearn(val as number)} marks={[{value: 0, label: '0'}, {value: 100, label: '100'}]} step={1} valueLabelDisplay="auto" />
        <Divider sx={{my: 2}} />
        <Typography variant="h6">Self Sufficient Learning</Typography>
        <Typography variant="subtitle2">
          Self-sufficient people have the ability and the desire to live life on their own terms, to determine their own course, to make their own decisions, and not have their life choices made by others.
          <br />Are they able to choose what to learn or do they need someone to tell them
        </Typography>
        <Slider value={selfSufficientLearning} onChange={(_, val) => setSelfSufficientLearning(val as number)} marks={[{value: 0, label: '0'}, {value: 100, label: '100'}]} step={1} valueLabelDisplay="auto" />
        <Divider sx={{my: 2}} />
        <Typography variant="h6">Improving Capability</Typography>
        <Typography variant="subtitle2">
          Do they actively seek out improving existing capabilities at their place of work/home and/or create new capabilities.
        </Typography>
        <Slider value={improvingCapability} onChange={(_, val) => setImprovingCapability(val as number)} marks={[{value: 0, label: '0'}, {value: 100, label: '100'}]} step={1} valueLabelDisplay="auto" />
        <Divider sx={{my: 2}} />
        <Typography variant="h6">Innovative Thinking</Typography>
        <Typography variant="subtitle2">
          Do you think of them as someone who comes up with new problems or solutions that no one else may have thought of yet?
        </Typography>
        <Slider value={innovativeThinking} onChange={(_, val) => setInnovativeThinking(val as number)} marks={[{value: 0, label: '0'}, {value: 100, label: '100'}]} step={1} valueLabelDisplay="auto" />
        <Divider sx={{my: 2}} />
        <Typography variant="h6">Growth Mindset</Typography>
        <Typography variant="subtitle2">
          Are they able to identify areas of growth needed in yourself/others?
        </Typography>
        <Slider value={growthMindset} onChange={(_, val) => setGrowthMindset(val as number)} marks={[{value: 0, label: '0'}, {value: 100, label: '100'}]} step={1} valueLabelDisplay="auto" />
        <Divider sx={{my: 2}} />
        <Typography variant="h6">Awareness Of Self Efficacy</Typography>
        <Typography variant="subtitle2">
          Self-efficacy refers to an individual&apos;s belief in his or her capacity to execute behaviors necessary to produce specific performance attainments (Bandura, 1977, 1986, 1997). Self-efficacy reflects confidence in the ability to exert control over one&apos;s own motivation, behavior, and social environment.
        </Typography>
        <Slider value={awarenessOfSelfEfficacy} onChange={(_, val) => setAwarenessOfSelfEfficacy(val as number)} marks={[{value: 0, label: '0'}, {value: 100, label: '100'}]} step={1} valueLabelDisplay="auto" />
        <Divider sx={{my: 2}} />
        <Typography variant="h6">Applying What They Learn</Typography>
        <Typography variant="subtitle2">
          Do they apply in practicality something that they have just been taught, e.g. if they learned cloud concepts recently can they actually use cloud concepts or do they still reluct to.
        </Typography>
        <Slider value={applyingWhatTheyLearn} onChange={(_, val) => setApplyingWhatTheyLearn(val as number)} marks={[{value: 0, label: '0'}, {value: 100, label: '100'}]} step={1} valueLabelDisplay="auto" />
        <Divider sx={{my: 2}} />
        <Typography variant="h6">Adaptability</Typography>
        <Typography variant="subtitle2">
          Can they change how they act/behave and/or can they improve based on circumstances that are thrown at them.
        </Typography>
        <Slider value={adaptability} onChange={(_, val) => setAdaptability(val as number)} marks={[{value: 0, label: '0'}, {value: 100, label: '100'}]} step={1} valueLabelDisplay="auto" />
        <Divider sx={{my: 2}} />
        <Button disabled={selfAssessment.isLoading || externalAssessment.isLoading} sx={{mr: 2}} onClick={() => navigate('/')}>Cancel</Button>
        <Button disabled={selfAssessment.isLoading || externalAssessment.isLoading} onClick={() => doAssessmentCreation()}>Create Assessment</Button>
      </Box>
    </>
  );
}
