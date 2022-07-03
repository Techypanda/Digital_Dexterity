interface Assessment {
  WillingnessToLearn: number;
  SelfSufficientLearning: number;
  ImprovingCapability: number;
  InnovativeThinking: number;
  GrowthMindset: number;
  AwarenessOfSelfEfficacy: number;
  ApplyingWhatTheyLearn: number;
  Adaptability: number;
}
interface AssessmentsResponse {
  externalAssessments?: Assessment[],
  selfAssessment?: Assessment,
}