import {
  Chart as ChartJS,
  RadialLinearScale,
  PointElement,
  LineElement,
  Filler,
  Tooltip,
  Legend,
} from 'chart.js';
import { Radar } from "react-chartjs-2";

ChartJS.register(
  RadialLinearScale,
  PointElement,
  LineElement,
  Filler,
  Tooltip,
  Legend
);

export function DigitalDexterityGraph(props: { externalAssessments?: Assessment[], selfAssessment?: Assessment }) {
  const averageExternal = () => {
    const avg: Assessment = {
      WillingnessToLearn: 0,
      SelfSufficientLearning: 0,
      ImprovingCapability: 0,
      InnovativeThinking: 0,
      GrowthMindset: 0,
      AwarenessOfSelfEfficacy: 0,
      ApplyingWhatTheyLearn: 0,
      Adaptability: 0
    };
    if (props.externalAssessments) {
      for (const assess of props.externalAssessments) {
        avg.Adaptability += assess.Adaptability;
        avg.ApplyingWhatTheyLearn += assess.ApplyingWhatTheyLearn;
        avg.AwarenessOfSelfEfficacy += assess.AwarenessOfSelfEfficacy;
        avg.GrowthMindset += assess.GrowthMindset;
        avg.ImprovingCapability += assess.ImprovingCapability;
        avg.InnovativeThinking += assess.InnovativeThinking;
        avg.SelfSufficientLearning += assess.SelfSufficientLearning;
        avg.WillingnessToLearn += assess.WillingnessToLearn;
      }
      avg.Adaptability /= props.externalAssessments.length;
      avg.ApplyingWhatTheyLearn /= props.externalAssessments.length;
      avg.AwarenessOfSelfEfficacy /= props.externalAssessments.length;
      avg.GrowthMindset /= props.externalAssessments.length;
      avg.ImprovingCapability /= props.externalAssessments.length;
      avg.InnovativeThinking /= props.externalAssessments.length;
      avg.SelfSufficientLearning /= props.externalAssessments.length;
      avg.WillingnessToLearn /= props.externalAssessments.length;
    }
    return [avg.Adaptability, avg.ApplyingWhatTheyLearn, avg.AwarenessOfSelfEfficacy, avg.GrowthMindset, avg.ImprovingCapability, avg.InnovativeThinking, avg.SelfSufficientLearning, avg.WillingnessToLearn];
  }
  function getDatasets() {
    const datasets = [];
    if (props.selfAssessment) {
      const temp = props.selfAssessment;
      datasets.push({
        label: 'What You Assessed Your Digital Dexterity As',
        data: [temp.Adaptability, temp.ApplyingWhatTheyLearn, temp.AwarenessOfSelfEfficacy, temp.GrowthMindset, temp.ImprovingCapability, temp.InnovativeThinking, temp.SelfSufficientLearning, temp.WillingnessToLearn],
        backgroundColor: 'rgba(54, 162, 235, 0.2)',
        borderColor: 'rgb(54, 162, 235)',
        borderWidth: 1,
      })
    }
    if (props.externalAssessments) {
      datasets.push({
        label: 'What Other Users Assess Your Digital Dexterity As',
        data: averageExternal(),
        backgroundColor: 'rgba(255, 99, 132, 0.2)',
        borderColor: 'rgba(255, 99, 132, 1)',
        borderWidth: 1,
      })
    }
    return datasets;
  }
  return (
    <Radar
      data={{
        labels: ['Adaptability', 'Applying What They Learn', 'Awareness Of Self Efficacy', 'Growth Mindset', 'Improving Capability', 'Innovative Thinking', 'Self Sufficient Learning', 'Willingness To Learn'],
        datasets: getDatasets(),
      }}
      options={{
        scales: {
          r: {
            min: 0,
            max: 100,
          } 
        }
      }}
    />
  )
}