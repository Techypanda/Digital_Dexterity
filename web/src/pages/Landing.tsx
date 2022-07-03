import { Alert, Box, Button, CircularProgress, Dialog, DialogContent, DialogTitle, Typography } from "@mui/material";
import { CommonProps } from "@mui/material/OverridableComponent";
import { useAssessments } from "../api/assessments";
import { DigitalDexterityGraph } from "../components/landing/DigitalDexterityGraph";
import { useNavigate } from "react-router-dom";
import { useState } from "react";

export default function Landing(props: CommonProps) {
  const { data, isLoading, isFetching, isError } = useAssessments();
  const [showAssessmentDialog, setShowAssessmentDialog] = useState(false);
  const navigate = useNavigate();
  return (
    <>
      <Typography variant="h6" component="h2" align="center" mt={2}>My Digital Dexterity</Typography>
      <Box mb={3}>
        <Button variant="contained" sx={{ mr: 2 }} onClick={() => setShowAssessmentDialog(true)}>Assess</Button>
        <Button variant="contained" onClick={() => navigate("/explain")}>What Do These Labels Mean?</Button>
      </Box>
      <Dialog open={showAssessmentDialog} onClose={() => setShowAssessmentDialog(false)}>
        <DialogTitle>What Assessment Type Do You Want To Fill Out?</DialogTitle>
        <DialogContent>
          <Box display="flex" justifyContent="space-around">
            <Button onClick={() => navigate("/assess/self")}>Assess Myself</Button>
            <Button onClick={() => navigate("/assess")}>Assess Someone Else</Button>
          </Box>
        </DialogContent>
      </Dialog>
      {(isLoading || isFetching) ?
        <Box display="flex" justifyContent="center" mt={4}>
          <CircularProgress size="6rem" />
        </Box>
        :
        <>
          {isError ? <Alert severity="error">An Error Has Occured, Please try a relogin and then if it persists contact an admin</Alert> :
            <DigitalDexterityGraph externalAssessments={data!.externalAssessments} selfAssessment={data!.selfAssessment} />
          }
        </>
      }
    </>
  )
}