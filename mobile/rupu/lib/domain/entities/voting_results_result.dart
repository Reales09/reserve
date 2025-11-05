class VotingResultsResult {
  final bool success;
  final String? message;
  final List<VotingResult> data;

  const VotingResultsResult({
    required this.success,
    this.message,
    required this.data,
  });
}

class VotingResult {
  final int id;
  final int votingId;
  final int propertyUnitId;
  final int votingOptionId;
  final DateTime votedAt;
  final String ipAddress;
  final String userAgent;

  const VotingResult({
    required this.id,
    required this.votingId,
    required this.propertyUnitId,
    required this.votingOptionId,
    required this.votedAt,
    required this.ipAddress,
    required this.userAgent,
  });
}
