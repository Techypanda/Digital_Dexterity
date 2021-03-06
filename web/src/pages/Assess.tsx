import {Paper} from '@mui/material';
import {CommonProps} from '@mui/material/OverridableComponent';
import AssessmentForm from '../components/assessment/AssessmentForm';

export default function Assess(props: CommonProps) {
  return (
    <Paper sx={{my: 4, p: 3}}>
      <AssessmentForm />
    </Paper>
  );
}
