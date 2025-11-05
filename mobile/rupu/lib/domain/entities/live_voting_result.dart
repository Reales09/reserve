class LiveVotingResult {
  final int totalUnits;
  final List<LiveVotingUnit> units;

  const LiveVotingResult({
    required this.totalUnits,
    required this.units,
  });
}

class LiveVotingUnit {
  final int propertyUnitId;
  final String propertyUnitNumber;
  final double participationCoefficient;
  final int? residentId;
  final String? residentName;
  final bool hasVoted;
  final int? votingOptionId;
  final String? optionText;
  final String? optionCode;
  final String? optionColor;
  final DateTime? votedAt;

  const LiveVotingUnit({
    required this.propertyUnitId,
    required this.propertyUnitNumber,
    required this.participationCoefficient,
    this.residentId,
    this.residentName,
    required this.hasVoted,
    this.votingOptionId,
    this.optionText,
    this.optionCode,
    this.optionColor,
    this.votedAt,
  });
}
