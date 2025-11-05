class VotingOption {
  final int id;
  final int votingId;
  final String optionText;
  final String optionCode;
  final String color;
  final int displayOrder;
  final bool isActive;

  const VotingOption({
    required this.id,
    required this.votingId,
    required this.optionText,
    required this.optionCode,
    required this.color,
    required this.displayOrder,
    required this.isActive,
  });
}
